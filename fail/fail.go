package fail

import (
	"encoding/json"
	"unicode"
)

type Mistake struct {
	s string
}

type Mistakes map[string][]*Mistake

func (ms Mistakes) Error() string {
	return "validation errors"
}

func (ms Mistakes) MarshalJSON() ([]byte, error) {
	out := make(map[string][]interface{})
	for k, mm := range ms {
		for _, m := range mm {
			out[k] = append(out[k], m.Error())
		}
	}
	return json.Marshal(out)
}

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

func IsPrintable(v string) *Mistake {
	// only printables
	printable := true
	for _, ch := range v {
		printable = unicode.IsPrint(ch)
		if !printable {
			return NewMistake("unprintable characters")
		}
	}
	return nil
}
