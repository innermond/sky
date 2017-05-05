package mysql

import (
	"database/sql"

	"github.com/innermond/sky/sky"
)

var (
	_ sky.Session = &Session{}
)

type Session struct {
	db *sql.DB
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}
