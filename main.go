package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/innermond/sky/http"
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

	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	all := &http.AllServicesHandler{
		PersonHandler: http.NewPersonHandler(personService),
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
	c := make(chan bool, 1)
	<-c
	srv.Close()
}
