package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// all session code adapted from http://www.gorillatoolkit.org/pkg/sessions
//Global variables
var store = sessions.NewCookieStore([]byte("secret"))
var mongoConnection, err = newMongoConnection()
var currentUser = ""
var currentUserBlogs []Blog

// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
func indexRoute(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "views/index.html")
}

func main() {
	// adapted from https://www.reddit.com/r/golang/comments/2tp5ho/updated_my_ggap_stack_web_app_tutorial_slothful/
	router := initRouter()
	port := os.Getenv("PORT")
	//Serves the file on port 4000 if it is run locally, this allows the hosted page to work as well as the local build
	if port == "" {
		server := &http.Server{

			Addr:    ":4000",
			Handler: router,
		}
		server.ListenAndServe()
	} else {
		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}
		server.ListenAndServe()
	}
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

	// adapted from https://auth0.com/blog/authentication-in-golang/
	r.Handle("/register", errorHandler(Register)).Methods("POST")
	r.Handle("/login", errorHandler(loginHandler)).Methods("POST")
	r.Handle("/logout", errorHandler(logoutHandler)).Methods("POST")
	r.Handle("/user", errorHandler(createBlog)).Methods("POST")
	r.Handle("/blogs", errorHandler(getBlogs)).Methods("GET")
	r.Handle("/user", errorHandler(getUserBlogs)).Methods("GET")
	r.Handle("/user", errorHandler(deleteBlogPost)).Methods("DELETE") // yet to be implemented
	r.Handle("/user", errorHandler(updateBlogPost)).Methods("PUT")    // Yet to be implemented
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
	// Connect to the mongo database
	s, err := mgo.Dial("mongodb://test:test@ds035006.mlab.com:35006/heroku_lzbj5rj0")
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	return s, err
}

//Structs for mongo objects
type (
	User struct {
		Name      string
		Username  string
		Password  string
		Email     string
		Blogposts []string
	}
)


type (
	LoginCreds struct {
		Username string
		Password string
	}
)

type (
	Blog struct {
		UniqueId  string   `json:"uniqueid"`
		Title     string   `json:"title"`
		Body      []string `json:"body"`
		Author    string   `json:"author"`
		Likes     int      `json:"likes"`
		CreatedOn int      `json:"createOn"`
		Comments  []Comment
	}
)

type (
	Comment struct {
		CBlogID string `json:"cblogid"`
		CBody   string `json:"cbody"`
		CAuthor string `json:"cauthor"`
	}
)

//Handler functions
func Register(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var user User
	//Pull in the user data from the webpage
	err := decoder.Decode(&user)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	//Put the data from the webpage into local go variables
	u := user.Username
	p := user.Password
	e := user.Email
	n := user.Name

	a := User{Username: u, Password: p, Email: e, Name: n}
	if a.Username != "" || a.Password != "" || a.Email != "" || a.Name != "" {
		insert(a)
	}
	return err
}


// adapted from https://devcenter.heroku.com/articles/go-sessions
func loginHandler(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var login LoginCreds
	//Extract the login details from the webpage
	err := decoder.Decode(&login)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	// Call login validation to check whether the login details are in the database
	if err := loginValidation(login.Username, login.Password); err == nil {
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//save the login details in the local session
		session.Values["username"] = login.Username
		session.Values["password"] = login.Password
		session.Save(r, w)
		w.Header().Add("username", currentUser)
		w.Header().Add("password", login.Password)
	} else {
		return err
	}
	return err
}


func logoutHandler(w http.ResponseWriter, req *http.Request) error {
	session, err := store.Get(req, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	// Set session login values to nil
	session.Values["username"] = ""
	session.Values["password"] = ""
	if err := session.Save(req, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Header().Add("username", "")
	w.Header().Add("password", "")
	currentUser = ""
	return err
}

// Function not working on the angular side
func updateBlogPost(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var blog Blog
	jErr := decoder.Decode(&blog)
	if err != nil {
		return jErr
	}
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Blogs")
	err = c.Update(bson.M{"uniqueId": blog.UniqueId},
		bson.M{"title": blog.Title})
	if err != nil {
		return err
	}
	return nil
}

//Get functions
func getComments(blogID string) []Comment {
	commentData := mongoConnection.DB("heroku_lzbj5rj0").C("Comments")
	var resultComments []Comment
	err = commentData.Find(bson.M{"cblogid": blogID}).All(&resultComments)
	if err != nil {
		log.Fatal(err)
	}
	if commentData != nil {
		return resultComments
	} else {
		return nil
	}
}

func getBlogs(w http.ResponseWriter, r *http.Request) error {
	var results []Blog
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Blogs")
	mErr := c.Find(nil).All(&results)
	if err != nil {
		return mErr
	}
	//Return the comments for each blog one by One
	//Returning comments as a apart of the blog could not be implemented effectively
	for x := 0; x <= len(results)-1; x++ {
		results[x].Comments = getComments(results[x].UniqueId)
	}
	json.NewEncoder(w).Encode(results)
	return nil
}

func getUserBlogs(w http.ResponseWriter, r *http.Request) error {
	currentUserBlogs = nil
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	resultingBlogID := User{}
	//Return blog id's from the user document
	err = c.Find(bson.M{"username": currentUser}).Select(bson.M{"blogposts": 1, "_id": 0}).One(&resultingBlogID)
	if err != nil {
		return err
	}
	//Check if the user has any blog posts if true then return the array
	if resultingBlogID.Blogposts != nil {
		blogData := mongoConnection.DB("heroku_lzbj5rj0").C("Blogs")
		resultBlog := Blog{}
		//Return user blogs and append
		for i := 0; i <= len(resultingBlogID.Blogposts)-1; i++ {
			err = blogData.Find(bson.M{"uniqueid": resultingBlogID.Blogposts[i]}).One(&resultBlog)
			resultBlog.Comments = getComments(resultingBlogID.Blogposts[i])
			if err != nil {
				return err
			}
			currentUserBlogs = append(currentUserBlogs, resultBlog)
		}
	}
	json.NewEncoder(w).Encode(currentUserBlogs)
	return nil
}

func loginValidation(username string, password string) error {
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	result := User{}
	err = c.Find(bson.M{"username": username}).Select(bson.M{"username": 1, "password": 1, "_id": 0}).One(&result)
	if err != nil {
		return err
	}
	if result.Username == username && result.Password == password {
		fmt.Println("Connection succesful")
		currentUser = username
		return err
	} else {
		return err
	}
}

//Insert Functions

//adapted from https://stevenwhite.com/building-a-rest-service-with-golang-3/ used to make connection to mongoDB database
func insert(a User) {
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	err = c.Insert(&User{a.Name, a.Username, a.Password, a.Email, nil})
	if err != nil {
		log.Fatal(err)
	}
}

func insertComment(a Comment) {
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Comments")
	err = c.Insert(&Comment{a.CBlogID, a.CBody, a.CAuthor})
	if err != nil {
		log.Fatal(err)
	}
}

func createBlog(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var blog Blog
	err := decoder.Decode(&blog)
	if err != nil {
		return err
	}
	resultingBlogID := User{}
	u := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	//Update the users list of blog ID's for reference
	err = u.Find(bson.M{"username": currentUser}).Select(bson.M{"blogposts": 1}).One(&resultingBlogID)
	resultingBlogID.Blogposts = append(resultingBlogID.Blogposts, blog.UniqueId)
	err = u.Update(bson.M{"username": currentUser}, bson.M{"$set": bson.M{"blogposts": resultingBlogID.Blogposts}})

	c := mongoConnection.DB("heroku_lzbj5rj0").C("Blogs")
	//Insert new blog into the mongo database
	err = c.Insert(&Blog{blog.UniqueId, blog.Title, blog.Body, blog.Author, blog.Likes, blog.CreatedOn, blog.Comments})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}


//Delete Functions
func deleteBlogPost(w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	var blog Blog
	err := decoder.Decode(&blog)

	fmt.Println("Remove actual blog")
	c := mongoConnection.DB("heroku_lzbj5rj0").C("Blogs")
	//Remove the blog from the collection
	err = c.Remove(bson.M{"uniqueid": blog.UniqueId})
	if err != nil {
		return err
	}

	//Remove the blog ID from the user ID array
	var tempBlogArray []string
	resultingBlogID := User{}
	u := mongoConnection.DB("heroku_lzbj5rj0").C("Users")
	err = u.Find(bson.M{"username": currentUser}).Select(bson.M{"blogposts": 1}).One(&resultingBlogID)
	if err != nil {
		return err
	}
	for x := 0; x <= len(resultingBlogID.Blogposts)-1; x++ {
		if resultingBlogID.Blogposts[x] != blog.UniqueId {
			tempBlogArray = append(tempBlogArray, (resultingBlogID.Blogposts[x]))
		}
	}
	err = u.Update(bson.M{"username": currentUser}, bson.M{"$set": bson.M{"blogposts": tempBlogArray}})
	if err != nil {
		return err
	}

	return nil
}

// adapted from https://github.com/campoy/todo/blob/master/server/server.go
// badRequest is handled by setting the status code in the reply to StatusBadRequest.
type badRequest struct{ error }

// notFound is handled by setting the status code in the reply to StatusNotFound.
type notFound struct{ error }

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "oops", http.StatusInternalServerError)
		}
	}
}
