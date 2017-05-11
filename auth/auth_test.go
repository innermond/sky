package auth

import (
	"crypto/rsa"
	"fmt"
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

func TestAuthenticator_createAuthToken_getClaims(t *testing.T) {
	verify, err := PublicKey(pubpath)
	if err != nil {
		t.Fatal(err)
	}

	sign, err := PrivateKey(privpath)
	if err != nil {
		t.Fatal(err)
	}

	a := NewAuthenticator(verify)
	tk := NewTokenCreator(sign)
	t.Logf("%v %v", a, tk)

	uses := []int{1, 2, 3, 4, 5}
	for _, i := range uses {
		t.Run(string(i), func(t *testing.T) {
			claims := &AuthClaims{
				&jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
				},
				UserInfo{i},
			}
			s, err := tk.CreateAuthToken(claims)
			if err != nil {
				t.Fatal(err)
			}
			t.Log(s)
			decoded, err := a.GetClaims(s)
			if err != nil {
				t.Fatal(err)
			}
			if decoded.UserInfo.Id != i {
				t.Errorf("expected %d got %d", i, decoded.UserInfo.Id)
			}
			t.Log(decoded)
		})
	}
}
func TestAuthenticator_authenticateTokenValidationError(t *testing.T) {
	verify, err := PublicKey(pubpath)
	if err != nil {
		t.Fatal(err)
	}

	a := NewAuthenticator(verify)
	t.Logf("%v", a)

	uses := []string{
		// expired
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ0ODE5NzcsIklkIjo0fQ.wnQCbTG0H5hEQ7cv5U-z-3g8Hgdu59x-CtSASmNXTrs5X6l4heqlIgLTXXke3djh0-HPhCM6ZvZR-tUBwgtmlor6H8txID8md6Yofo_nzeO_eGPYcXc2huVTN1Dpi7FIdji9t3kk_F_KNIJon8taMHKOY62MvtQwDhYG_3lPkNQ",
		// invalid
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ0ODE5NzcsIklkIjo0fQ.wnQCbTG0H5hEQ7cv5U-z-3g8Hgdu59x-CtSASmNXTrs5X6l4heqlIgLTXXke3djh0-HPhCM6ZvZR-tUBwgtmlor6H8txID8md6Yofo_nzeO_eGPYcXc2huVTN1Dpi7FIdji9t3kk_F_KNIJon8taMHKOY62MvtQwDhYG_3lP",
		"aaa",
	}
	for i, s := range uses {
		t.Run(string(i), func(t *testing.T) {
			err := a.Authenticate(s)
			if _, ok := err.(*jwt.ValidationError); !ok || err == nil {
				t.Fatal(s)
			}
			t.Log(fmt.Sprintf("%T %v", err, err))
		})
	}
}
