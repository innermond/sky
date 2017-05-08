package fail

type Mistake struct {
	s string
}

type Mistakes map[string][]*Mistake

func (e Mistake) Error() string {
	return e.s
}

func NewMistake(s string) *Mistake {
	return &Mistake{s}
}

type Failer interface {
	Fail() bool
	Err() Mistakes
}
