package mysql

import (
	"database/sql"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/innermond/sky/auth"
	"github.com/innermond/sky/sky"
)

var (
	_ sky.TokenService = &TokenService{}
)

type TokenService struct {
	session      *Session
	tokenCreator *auth.Creator
}

func NewTokenService(s *Session, c *auth.Creator) *TokenService {
	return &TokenService{s, c}
}

func (s *TokenService) Create(k sky.ApiKey) (string, error) {
	var (
		uid int
	)
	//q := `select id, if(password=sha2(concat(?, api_key), 256), 1, 0) as ok from users where username=? and api_key=? limit 1;`
	q := `select id from users where api_key=? limit 1;`
	key := string(k)
	err := s.session.db.QueryRow(q, key).Scan(&uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("userless apikey")
		}
		return "", err
	}
	claims := &auth.AuthClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		auth.UserInfo{uid},
	}

	tokstr, err := s.tokenCreator.CreateAuthToken(claims)
	if err != nil {
		return "", err
	}
	return tokstr, nil
}
