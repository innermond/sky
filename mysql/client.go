package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Client struct {
	db *sql.DB
	tx *sql.Tx
}

func (c *Client) Open(dns string) error {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return err
	}
	c.db = db

	if err = c.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Begin() error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	c.tx = tx
}
