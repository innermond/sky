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

type TokenService interface {
	Create(ApiKey) (string, error)
}

type ApiKey string
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session interface {
}

type Authenticator interface {
	Authenticate(string) error
}
