package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mgo.v2"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/templates/index.html"))
}

func main() {
	//m := macaron.Classic()
	//insert()
	http.HandleFunc("/", foo)
	http.HandleFunc("/css/", serveResource)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":4000", nil)
	//m.Run()
}

type (
	User struct {
		Username string
		Password string
	}
)

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

func foo(w http.ResponseWriter, req *http.Request) {

	u := req.FormValue("username")
	p := req.FormValue("password")
	err := tpl.ExecuteTemplate(w, "index.html", User{u, p})
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}

	a := User{Username: u, Password: p}
	if a.Username != "" || a.Password != "" {
		insert(a)
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
	err = c.Insert(&User{a.Username, a.Password})
	if err != nil {
		log.Fatal(err)
	}

	return s
}
