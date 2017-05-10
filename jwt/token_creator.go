package jwt

import (
	"crypto/rsa"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

type Creator struct {
	signKey *rsa.PrivateKey
}

func NewTokenCreator(piv *rsa.PrivateKey) *Creator {
	return &Creator{piv}
}

func (a *Creator) CreateAuthToken(claims *AuthClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, err := t.SignedString(a.signKey)
	log.Println("token", claims, s)
	if err != nil {
		return "", err
	}
	return s, nil
}
