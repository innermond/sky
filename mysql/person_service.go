package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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
