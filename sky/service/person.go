package service

type PersonService struct{}

func (service *PersonService) Get(pid sky.PersonID) (sky.Person, error) {
	return sky.Person{1, "test name"}, nil
}
