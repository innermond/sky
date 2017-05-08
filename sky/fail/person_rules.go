package fail

import (
	"github.com/innermond/sky/sky"
)

type PersonRules struct {
	sky.Person
	mistakes Mistakes
}

func (r *PersonRules) LongnameOk() *Mistake {
	v := r.Longname
	// required
	if v == "" {
		return NewMistake("required")
	}
	// size
	l := len(v)
	if l < 4 || l > 10 {
		return NewMistake("unexpected length")
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
		r.mistakes["longname"] = append(r.mistakes["longname"], merr)
		failed = true
	}

	return failed
}

func (r *PersonRules) Err() Mistakes {
	return r.mistakes
}
