package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// all session code adapted from http://www.gorillatoolkit.org/pkg/sessions
var store = sessions.NewCookieStore([]byte("secret"))
var mongoConnection, err = newMongoConnection()

// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/

func indexRoute(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func main() {
	// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
	router := initRouter()

	server := &http.Server{
		Addr:    ":4000",
		Handler: router,
	}

	fmt.Println("Starting server")

	server.ListenAndServe()
}

// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
func FileServerRouteG(m *mux.Router, path, dir string) {
	m.PathPrefix(path).Handler(
		http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}

// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
func AddStaticRoutes(m *mux.Router, pathsAndDirs ...string) {
	for i := 0; i < len(pathsAndDirs)-1; i += 2 {
		FileServerRouteG(m, pathsAndDirs[i], pathsAndDirs[i+1])
	}
}

// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
func initRouter() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/register", http.HandlerFunc(RegisterHandler)).Methods("POST")
	//Add static routes for the public directory
	AddStaticRoutes(r, "/partials/", "public/partials",
		"/scripts/", "public/scripts", "/styles/", "public/styles",
		"/images/", "public/images")

	//Serve all other requests with index.html, and ultimately the front-end
	//Angular.js app.
	r.PathPrefix("/").HandlerFunc(indexRoute)
	return r
}

func newMongoConnection() (*mgo.Session, error) {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://test:test@ds035006.mlab.com:35006/heroku_lzbj5rj0")
	fmt.Println("Mongo connected")
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s, err
}

type (
	User struct {
		Name     string
		Username string
		Password string
		Email    string
	}
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	u := req.FormValue("username")
	p := req.FormValue("password")
	e := req.FormValue("email")
	n := req.FormValue("name")
	//fmt.Println(u)
	session, err := store.Get(req, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.Values["username"] = u
	session.Values["password"] = p
	session.Save(req, w)

	a := User{Username: u, Password: p, Email: e, Name: n}
	if a.Username != "" || a.Password != "" || a.Email != "" || a.Name != "" {
		insert(a)
	}

	http.Redirect(w, req, "/", 302)
}

// adapted from https://devcenter.heroku.com/articles/go-sessions
func loginHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Login handler started")
	username := req.FormValue("username")
	password := req.FormValue("password")
	/* WE NEED TO ADD MONGO CHECKING HERE AS WELL */
	//if err := session.DB(authDB).Login(user, pass); err == nil {
	if loginValidation(username, password) {
		session, err := store.Get(req, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		session.Values["username"] = username
		session.Values["password"] = password
		session.Save(req, w)
		//}
		http.Redirect(w, req, "/", 302)
	} else {
		fmt.Println("Invalid login")
		// TODO: notify user of invalid username password
	}
}

func loginValidation(username string, password string) bool {
	fmt.Println("Login validation started")
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	result := User{}
	err = c.Find(bson.M{"username": username}).Select(bson.M{"username": 1, "password": 1, "_id": 0}).One(&result)
	if err != nil {
		// TODO: This exits the cript if the query fails to find the user, needs to be changed
		log.Fatal(err)
	}
	if result.Username == username && result.Password == password {
		fmt.Println("Connection succesful")
		return true
	} else {
		return false
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	session, err := store.Get(req, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["username"] = ""
	if err := session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/", 302)
}

//adapted from https://stevenwhite.com/building-a-rest-service-with-golang-3/ used to make connection to mongoDB database
func insert(a User) {
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	err = c.Insert(&User{a.Name, a.Username, a.Password, a.Email})
	if err != nil {
		log.Fatal(err)
	}
}
