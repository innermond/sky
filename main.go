package main

import (
	"fmt"
	"net/http"
)

// we are on earth

func main() {
	// create a handler
	indexHandler := IndexHandler{}
	// create a tcp server
	http.Handle("/", indexHandler)
	http.ListenAndServe(":3000", nil)
}

type IndexHandler struct{}

func (hnd IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// echo back url parameters
	s := fmt.Sprintf("%v", r.URL.Query())
	if r.Method == "POST" {
		r.ParseForm()
		s = fmt.Sprintf("%v", r.Form)
	}
	w.Write([]byte(s))
}
