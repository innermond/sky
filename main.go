package main

import (
	"log"
	"os"

	"github.com/innermond/sky/http"
)

// we are on earth

func main() {
	//logger
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	all := &http.AllServicesHandler{
		PersonHandler: http.NewPersonHandler(),
	}
	// create server
	srv := &http.IndexServer{
		Addr:    ":3000",
		Handler: all,
	}

	err := srv.Open()
	if err != nil {
		srv.Close()
	}
	c := make(chan bool, 1)
	<-c
	srv.Close()
}
