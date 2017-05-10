package jwt

import (
	"crypto/rsa"
	"io/ioutil"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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

func CreateAuthToken(claims *AuthClaims, signKey *rsa.PrivateKey) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return s, nil
}

func GetClaims(tokstr string, verifyKey *rsa.PublicKey) (*AuthClaims, error) {
	var claims AuthClaims
	tok, err := jwt.ParseWithClaims(tokstr, &claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}
	decoded := tok.Claims.(*AuthClaims)
	return decoded, nil
}

func TestAuthenticator_new(t *testing.T) {
	verify, err := PublicKey(pubpath)
	if err != nil {
		t.Fatal(err)
	}

	a := NewAuthenticator(verify)
	t.Logf("%v", a)

	sign, err := PrivateKey(privpath)
	if err != nil {
		t.Fatal(err)
	}

	uses := []int{1, 2, 3, 4, 5}
	for _, i := range uses {
		t.Run(string(i), func(t *testing.T) {
			claims := &AuthClaims{
				&jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
				},
				UserInfo{i},
			}
			s, err := CreateAuthToken(claims, sign)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(s)

			decoded, err := a.GetClaims(s)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(decoded)
		})
	}
}
