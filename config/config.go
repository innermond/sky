package config

import (
	"crypto/rsa"
	"database/sql"
	"io/ioutil"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

func init() {
	hm, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	root = "/home/" + hm + "/Projects/gollum/src/github.com/innermond/sky/"
}

var (
	pkgErr error
	root   string
	dns    = "root:M0b1d1c3@tcp(localhost:3306)/printoo"

	pubpath  string = root + "var/app.rsa.pub"
	privpath string = root + "var/app.rsa"
)

func Err() error {
	return pkgErr
}

func PublicKey() *rsa.PublicKey {
	if pkgErr != nil {
		return nil
	}
	pubBytes, err := ioutil.ReadFile(pubpath)
	if err != nil {
		pkgErr = err
		return nil
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		pkgErr = err
		return nil
	}
	return publicKey
}

func PrivateKey() *rsa.PrivateKey {
	if pkgErr != nil {
		return nil
	}
	privBytes, err := ioutil.ReadFile(privpath)
	if err != nil {
		pkgErr = err
		return nil
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		pkgErr = err
		return nil
	}
	return privateKey
}

func DB() *sql.DB {
	if pkgErr != nil {
		return nil
	}
	// database
	db, err := sql.Open("mysql", dns)
	if err != nil {
		pkgErr = err
		return nil
	}
	if err = db.Ping(); err != nil {
		pkgErr = err
		return nil
	}
	return db
}
