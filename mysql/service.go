package mysql

import "github.com/innermond/sky/sky"

var _ sky.PersonService = &PersonService{}

type PersonService struct {
}

func (p *sky.PersonService) Get(pid sky.PersonID) sky.Person {

}
