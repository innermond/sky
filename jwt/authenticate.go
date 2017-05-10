package jwt

import (
	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
)

type UserInfo struct {
	Id int
}

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
	var err error
	if a.token == nil {
		err = a.Authenticate(tokstr)
		if err != nil {
			return nil, err
		}
	}
	decoded := a.token.Claims.(*AuthClaims)
	return decoded, nil
}
