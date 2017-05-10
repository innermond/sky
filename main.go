package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/innermond/sky/config"
	"github.com/innermond/sky/http"
	myjwt "github.com/innermond/sky/jwt"
	"github.com/innermond/sky/mysql"
)

// we are on earth

func main() {
	//logger
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	dns := "root:M0b1d1c3@tcp(localhost:3306)/printoo"
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// authenticator
	verify, err := config.PublicKey()
	if err != nil {
		panic(err)
	}
	sign, err := config.PrivateKey()
	if err != nil {
		panic(err)
	}
	a := myjwt.NewAuthenticator(verify)
	c := myjwt.NewTokenCreator(sign)
	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	tokenService := mysql.NewTokenService(s, c)
	all := &http.AllServicesHandler{
		Auth:          a,
		PersonHandler: http.NewPersonHandler(personService),
		TokenHandler:  http.NewTokenHandler(tokenService),
	}
	// here wrap all handler with middlewares
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
