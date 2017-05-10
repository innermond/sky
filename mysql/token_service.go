package mysql

import (
	"database/sql"
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	myjwt "github.com/innermond/sky/jwt"
	"github.com/innermond/sky/sky"
)

var (
	_ sky.TokenService = &TokenService{}
)

type TokenService struct {
	session      *Session
	tokenCreator *myjwt.Creator
}

func NewTokenService(s *Session, c *myjwt.Creator) *TokenService {
	return &TokenService{s, c}
}

func (s *TokenService) Create(c sky.Credentials) (string, error) {
	var (
		uid int
		ok  bool
	)
	q := `select id, if(password=sha2(concat(?, salt), 256), 1, 0) as ok from users where username=? and salt=? limit 1;`
	err := s.session.db.QueryRow(q, c.Password, c.Username, c.ApiKey).Scan(&uid, &ok)
	log.Println("token service", err)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("cannot create token")
		}
		return "", err
	}
	claims := &myjwt.AuthClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		myjwt.UserInfo{uid},
	}

	tokstr, err := s.tokenCreator.CreateAuthToken(claims)
	log.Println("token service", uid, ok, tokstr)
	if err != nil {
		return "", err
	}
	return tokstr, nil
}
