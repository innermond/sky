package main

import (
	"crypto/rsa"
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/innermond/sky/http"
	myjwt "github.com/innermond/sky/jwt"
	"github.com/innermond/sky/mysql"
)

// we are on earth
const root string = "/home/userescu/Projects/gollum/src/github.com/innermond/sky/"

var (
	pubpath  string = root + "var/app.rsa.pub"
	privpath string = root + "var/app.rsa"
)

func PublicKey(pubPath string) (*rsa.PublicKey, error) {
	pubBytes, err := ioutil.ReadFile(pubPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

func PrivateKey(privPath string) (*rsa.PrivateKey, error) {
	privBytes, err := ioutil.ReadFile(privPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}

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
	verify, err := PublicKey(pubpath)
	if err != nil {
		panic(err)
	}
	a := myjwt.NewAuthenticator(verify)
	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	all := &http.AllServicesHandler{
		Auth:          a,
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
