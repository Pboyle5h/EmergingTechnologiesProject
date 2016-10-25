package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/templates/index.html"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", display)
	r.HandleFunc("/Register", Register)
	//r.HandleFunc("/css/", serveResource)
	http.Handle("/", r)
	http.HandleFunc("/css/", serveResource)
	http.ListenAndServe(":4000", nil)

}

func display(w http.ResponseWriter, req *http.Request) {
	err := tpl.Execute(w, "index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}

}

type (
	User struct {
		Name     string
		Username string
		Password string
		Email    string
	}
)

func Register(w http.ResponseWriter, req *http.Request) {

	u := req.FormValue("username")
	p := req.FormValue("password")
	e := req.FormValue("email")
	n := req.FormValue("name")
	err := tpl.ExecuteTemplate(w, "index.html", User{u, p, e, n})
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}

	a := User{Username: u, Password: p, Email: e, Name: n}
	if a.Username != "" || a.Password != "" || a.Email != "" || a.Name != "" {
		insert(a)
	}
}

//http://stackoverflow.com/questions/36323232/golang-css-files-are-being-sent-with-content-type-text-plain
func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

//adapted from https://stevenwhite.com/building-a-rest-service-with-golang-3/ used to make connection to mongoDB database
func insert(a User) *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://test:test@ds035006.mlab.com:35006/heroku_lzbj5rj0")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	c := s.DB("heroku_lzbj5rj0").C("Users")
	err = c.Insert(&User{a.Name, a.Username, a.Password, a.Email})
	if err != nil {
		log.Fatal(err)
	}

	return s
}
