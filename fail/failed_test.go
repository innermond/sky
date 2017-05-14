package fail

import (
	"strings"
	"testing"

	"github.com/innermond/sky/sky"
)

func TestPersonRulesFailure(t *testing.T) {
	uses := []struct {
		data     sky.Person
		expected bool
	}{
		// All this cases will fail
		// empty
		{sky.Person{Longname: ""}, true},
		// min length
		{sky.Person{Longname: "min"}, true},
		// max length
		{sky.Person{Longname: "overflow the maximum length"}, true},
		// spaces
		{sky.Person{Longname: " "}, true},
		// numbers inside
		{sky.Person{Longname: "I contain an 1"}, true},
		// unprintables inside
		{sky.Person{Longname: `Break
		able`}, true},
		{sky.Person{Longname: "Break\nable"}, true},
		{sky.Person{Longname: "Break\table"}, true},
		// non-letters inside
		{sky.Person{Longname: "Break able!!"}, true},
		// utf8 char counted as 1 char in calculated length
		{sky.Person{Longname: strings.Repeat("é", 11)}, true},
		{sky.Person{Longname: strings.Repeat("и", 11)}, true},
		// must not fail
		{sky.Person{Longname: "éc cпутник"}, false},
	}
	for _, uc := range uses {
		t.Run(uc.data.Longname, func(t *testing.T) {
			r := &PersonRules{uc.data, nil}
			if r.Fail() != uc.expected {
				t.Error(uc.data.Longname)
			} else {
				merr := r.Err()
				t.Log(merr)
			}
		})
	}
}
