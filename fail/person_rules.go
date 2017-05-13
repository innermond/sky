package fail

import (
	"regexp"
	"strings"

	"github.com/innermond/sky/sky"
)

type PersonRules struct {
	sky.Person
	mistakes Mistakes
}

func NewPersonRules(p sky.Person) *PersonRules {
	return &PersonRules{p, Mistakes{}}
}

func (r *PersonRules) LongnameOk() *Mistake {
	//TODO Move sanitising outside
	v := strings.TrimSpace(r.Person.Longname)
	// only printables
	if merr := IsPrintable(v); merr != nil {
		return merr
	}
	// required
	if v == "" {
		return NewMistake("required")
	}
	// size
	l := len([]rune(v))
	if l < 4 || l > 10 {
		return NewMistake("unexpected length")
	}
	// utf-8 letters
	fit, err := regexp.MatchString("[\\p{L}\\-]+", v)
	if !fit || err != nil {
		return NewMistake("unacceptable characters")
	}
	return nil
}

func (r *PersonRules) Fail() bool {
	var (
		merr   *Mistake
		failed bool
	)
	// reset all mistakes
	r.mistakes = Mistakes{}
	merr = r.LongnameOk()
	// check fields are ok
	if merr != nil {
		r.addMistake("longname", merr)
		failed = true
	}

	return failed
}

func (r *PersonRules) addMistake(key string, merr *Mistake) {
	r.mistakes[key] = append(r.mistakes[key], merr)
}

func (r *PersonRules) Err() Mistakes {
	return r.mistakes
}
