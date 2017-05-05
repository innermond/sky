package sky

type PersonID int
type Person struct {
	ID       PersonID `json:"id,omitempty"`
	Longname string   `json:"longname"`
}

type PersonService interface {
	Get(PersonID) (*Person, error)
	Create(Person) (PersonID, error)
}

type Session interface {
}
