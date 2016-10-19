package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/sessions"
)

// func getCode(data string) string {
// 	h := hmac.New(sha256.New, []byte("ourkey"))
// 	io.WriteString(h, data)
// 	return fmt.Sprintf("%x", h.Sum(nil))
// }

var store = sessions.NewCookieStore([]byte("secret password"))

func main() {

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {

		session, _ := store.Get(req, "session")
		if req.FormValue("email") != "" {
			// usually check password here

			session.Values["email"] = req.FormValue("email")
		}
		session.Save(req, res)
		// cookie, err := req.Cookie("session-id")
		// if err != nil { // cookie not set
		// 	//id, _ := uuid.NewV4()
		// 	cookie = &http.Cookie{
		// 		Name:  "session-id",
		// 		Value: req.FormValue("email"),
		// 	}
		// }
		// if req.FormValue("email") != "" {
		// 	cookie.Value = req.FormValue("email")
		// }
		//
		// code := getCode(cookie.Value)
		// cookie.Value = code + "|" + cookie.Value
		//
		// http.SetCookie(res, cookie)

		io.WriteString(res, `<!DOCTYPE html>
			<html>
				<body>
					<form>
						`+fmt.Sprint(session.Values["email"])+
			`<input type="email" name="email">
						 <input type="password" name="password">
						<input type="submit">
					</form>
				</body>
			</html>`)
	})

	http.ListenAndServe(":9000", nil)
}
