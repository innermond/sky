package main

import (
	"fmt"
	"net"
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

	// create server
	srv := &IndexServer{
		Addr:    ":3000",
		Handler: mux,
	}

	err := srv.Open()
	if err != nil {
		srv.Close()
	}
	c := make(chan bool, 1)
	<-c
	srv.Close()
}

type IndexServer struct {
	ln      net.Listener
	Addr    string
	Handler http.Handler
}

func (s *IndexServer) Open() error {
	ln, err := net.Listen("tcp", ":3000")
	s.ln = ln
	if err != nil {
		return err
	}
	go func() { http.Serve(s.ln, s.Handler) }()
	return nil
}

func (s *IndexServer) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
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
