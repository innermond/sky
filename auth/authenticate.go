package auth

import (
	"crypto/rsa"
	"errors"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserInfo struct {
	Id int
}

var ErrDecodingClaims = errors.New("error decoding claims")

type AuthClaims struct {
	*jwt.StandardClaims
	UserInfo
}

type Authenticator struct {
	verifyKey *rsa.PublicKey
	claims    *AuthClaims

	token *jwt.Token
}

func NewAuthenticator(pub *rsa.PublicKey) *Authenticator {
	claims := &AuthClaims{}
	return &Authenticator{pub, claims, nil}
}

func (a *Authenticator) Authenticate(tokstr string) error {
	tok, err := jwt.ParseWithClaims(
		tokstr,
		a.claims,
		func(token *jwt.Token) (interface{}, error) {
			return a.verifyKey, nil
		})

	a.token = tok

	if err != nil {
		return err
	}
	return nil
}

func (a *Authenticator) GetClaims(tokstr string) (*AuthClaims, error) {
	err := a.Authenticate(tokstr)
	if err != nil {
		return nil, err
	}
	decoded, ok := a.token.Claims.(*AuthClaims)
	if !ok {
		return decoded, ErrDecodingClaims
	}
	return decoded, nil
}
