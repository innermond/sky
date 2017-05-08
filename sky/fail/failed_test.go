package fail

import (
	"testing"

	"github.com/innermond/sky/sky"
)

func TestPersonRulesFailure(t *testing.T) {
	uses := []struct {
		data     sky.Person
		expected bool
	}{
		// all this cases will fail
		{sky.Person{Longname: ""}, true},
		{sky.Person{Longname: "min"}, true},
		{sky.Person{Longname: "overflow the maximum langth"}, true},
		{sky.Person{Longname: " "}, true},
		{sky.Person{Longname: "I contain an 1"}, true},
		{sky.Person{Longname: `Break
		able`}, true},
		{sky.Person{Longname: "Break\nable"}, true},
		{sky.Person{Longname: "Break\table"}, true},
	}
	for _, uc := range uses {
		t.Run(uc.data.Longname, func(t *testing.T) {
			r := &PersonRules{uc.data, nil}
			t.Log(r.Err())
			if r.Fail() != uc.expected {
				t.Error(uc.data.Longname)
			}
		})
	}
}
