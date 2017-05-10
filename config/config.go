package config

import (
	"crypto/rsa"
	"io/ioutil"

	jwt "github.com/dgrijalva/jwt-go"
)

const root string = "/home/gabriel/Projects/gollum/src/github.com/innermond/sky/"

var (
	pubpath  string = root + "var/app.rsa.pub"
	privpath string = root + "var/app.rsa"
)

func PublicKey() (*rsa.PublicKey, error) {
	pubBytes, err := ioutil.ReadFile(pubpath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, err
}

func PrivateKey() (*rsa.PrivateKey, error) {
	privBytes, err := ioutil.ReadFile(privpath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, err
}
