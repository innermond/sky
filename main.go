package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/innermond/sky/auth"
	"github.com/innermond/sky/config"
	"github.com/innermond/sky/http"
	"github.com/innermond/sky/mysql"
)

// we are on earth

func main() {
	//logger
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	// database
	db := config.DB()
	verify := config.PublicKey()
	sign := config.PrivateKey()
	err := config.Err()
	if err != nil {
		panic(err)
	}

	a := auth.NewAuthenticator(verify)
	c := auth.NewTokenCreator(sign)

	// session
	s := mysql.NewSession(db)

	// services
	personService := mysql.NewPersonService(s)
	tokenService := mysql.NewTokenService(s, c)

	all := &http.AllServicesHandler{
		AllRoutesAuth: a,
		PersonHandler: http.NewPersonHandler(personService),
		TokenHandler:  http.NewTokenHandler(tokenService),
	}

	// create server
	srv := &http.IndexServer{
		Addr:    ":3000",
		Handler: all,
	}

	err = srv.Open()
	if err != nil {
		srv.Close()
	}
	log.Println("serving...")
	loop := make(chan bool, 1)
	<-loop
	srv.Close()
}
