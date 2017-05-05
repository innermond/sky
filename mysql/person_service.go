package mysql

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/innermond/sky/sky"
)

var (
	_ sky.PersonService = &PersonService{}
)

type PersonService struct {
	session *Session
}

func NewPersonService(s *Session) *PersonService {
	return &PersonService{session: s}
}

func (s *PersonService) Get(pid sky.PersonID) (*sky.Person, error) {
	id := int(pid)
	q := `select
	id,
	longname/*,
	phone,
	email,
	(is_male=true),
	address,
	is_client,
	is_contractor*/
	from persons
	where id=?`

	var p sky.Person
	err := s.session.db.QueryRow(q, id).Scan(&p.ID, &p.Longname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (s *PersonService) Create(p sky.Person) (sky.PersonID, error) {
	q := `insert into persons (longname) values (?)`
	stm, err := s.session.db.Prepare(q)
	if err != nil {
		return 0, err
	}
	log.Println(p)
	res, err := stm.Exec(p.Longname)
	log.Println(err)
	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1364 {
				err = sky.ErrPersonValid
			}
		}
		return 0, err
	}
	lid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return sky.PersonID(lid), nil
}
