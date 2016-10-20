package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

// handle templates
func handleSomeRoute(res http.ResponseWriter, req *http.Request) {
	// step 1 parse template
	tmpl, err := template.ParseFiles("assets/templates/index.tmpl")
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	// execuate template
	tmpl.ExecuteTemplate(res, "assets/templates/index.tmpl", MyModel{
		SomeFields: 123,
		Values:     []int{1, 2, 3, 4, 5},
		Data:       "Some Value",
	})
}

// types of stuff on the page
type MyModel struct {
	SomeFields int
	Values     []int
	Data       string
}

// sessions
var store = sessions.NewCookieStore([]byte("secretcode"))

func handleSomeSessionRoute(res http.ResponseWriter, req *http.Request) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(req, "session-name")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	value := session.Values["test"] // get a Value
	// get value as a string
	str, _ := session.Values["test"].(string)
	fmt.Println(value, str)

	session.Values["test"] = 43 // set a Value

	delete(session.Values, "test") // delete a value
	session.Save(req, res)
}

// login example
func loginPage(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "session")
	if req.Method == "POST" {
		email := req.FormValue("email")
		password := req.FormValue("password")
		if email == "whatever" && password == "some-password" {
			session.Values["logged_in"] = "YES"
		} else {
			http.Error(res, "invalid details", 401)
		}
	}
}

func main() {

	http.HandleFunc("/test", handleSomeRoute)
	http.HandleFunc("/session", handleSomeSessionRoute)
	http.ListenAndServe(":8080", nil)
}
