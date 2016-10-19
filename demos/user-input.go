package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		key := "stuff"
		//val := req.URL.Query().Get(key)
		val := req.FormValue(key)
		//file, hdr, err := req.FormFile(key)
		//if err != nil {
		//	panic(err)
		//}
		//defer file.Close()
		//bs, _ := ioutil.ReadAll(file)
		//fmt.Println(string(bs))
		//fmt.Println(file, hdr, err)
		fmt.Println("Value: " + val)
		//io.WriteString(res, "Do my search for : "+val)
		res.Header().Set("Content-Type", "text/html")
		//io.WriteString(res, `<form method="POST"><input type="text" name="q"><input type="submit"></form>`)
		//io.WriteString(res, `<form method="POST"><input type="checkbox" name="q"><input type="submit"></form>`)
		io.WriteString(res, `<form method="POST" enctype="multipart/form-data"><input type="file" name="stuff"><input type="submit"></form>`)
	})
	http.ListenAndServe(":9000", nil)
}
