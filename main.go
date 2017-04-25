package main

import (
	"fmt"
	"net/http"
)

// we are on earth

func main() {
	//
	// create a handler
	indexHandler := IndexHandler{"earth is flat"}
	// create a tcp server
	mux := http.NewServeMux()
	mux.Handle("/", indexHandler)

	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	server.ListenAndServe()
}

// we can have more data inside, funcs, db handlers, encoders, ...
type IndexHandler struct {
	mySecret string
}

func (hnd IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// echo back url parameters
	s := fmt.Sprintf("%v", r.URL.Query())
	if r.Method == "POST" {
		r.ParseForm()
		s = fmt.Sprintf("%v", r.Form)
	}
	s = fmt.Sprintf("%v %q", s, hnd.mySecret)
	w.Write([]byte(s))
}
