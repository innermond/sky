package mock

type HappyAuthenticator struct{}

func NewHappyAuthenticator() *HappyAuthenticator {
	return &HappyAuthenticator{}
}

func (ha *HappyAuthenticator) Authenticate(tokstr string) error {
	return nil
}
