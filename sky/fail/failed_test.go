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
		{sky.Person{Longname: ""}, true},
		{sky.Person{Longname: "min"}, true},
		{sky.Person{Longname: "overflow the maximum langth"}, true},
	}
	for _, uc := range uses {
		t.Run(uc.data.Longname, func(t *testing.T) {
			r := &PersonRules{uc.data, nil}
			t.Log(r.Err())
			if r.Fail() != uc.expected {
				t.Error("fail")
			}
		})
	}
}
