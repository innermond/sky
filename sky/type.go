package sky

type PersonID int
type Person struct {
	ID       PersonID `json:"id,omitempty"`
	Longname string   `json:"longname"`
}

type PersonService interface {
	Get(PersonID) (*Person, error)
	Delete(PersonID) error
	Create(Person) (PersonID, error)
	Modify(Person) error
}

type Session interface {
}

type Authenticator interface {
	Authenticate() error
}
