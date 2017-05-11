package mysql

import (
	"strconv"
	"testing"

	"github.com/innermond/sky/config"
	"github.com/innermond/sky/sky"
)

func TestPersonService_get(t *testing.T) {
	db := config.DB()
	if err := config.Err(); err != nil {
		t.Fatal(err)
	}

	s := NewSession(db)
	ps := NewPersonService(s)

	userids := map[int]bool{
		1:    true,
		2:    true,
		3:    true,
		4:    true,
		5:    true,
		1000: true,
		2000: true,
	}
	for id, ok := range userids {
		t.Run(string(id), func(t *testing.T) {
			k := sky.PersonID(id)
			p, err := ps.Get(k)
			if err != nil && ok {
				t.Error(err)
			}
			if err == nil && !ok {
				t.Error(err)
			}
			if (id == 1000 || id == 2000) && !(p == nil && err == nil) {
				t.Errorf("p expected nil got %v and err expected nil got %v", p, err)
			}
			t.Log(ok, p)
		})
	}
}

func TestPersonService_create_delete(t *testing.T) {
	db := config.DB()
	if err := config.Err(); err != nil {
		t.Fatal(err)
	}

	s := NewSession(db)
	ps := NewPersonService(s)

	persons := []sky.Person{
		{Longname: "test1"},
		{Longname: "test2"},
	}
	for _, p := range persons {
		t.Run(p.Longname, func(t *testing.T) {
			pid, err := ps.Create(p)
			t.Log("personId", pid)
			if err != nil {
				t.Error(err)
			}
			err = ps.Delete(pid)
			if err != nil {
				t.Error("delete phase", err)
			}
		})
	}
}

func TestPersonService_delete(t *testing.T) {
	db := config.DB()
	if err := config.Err(); err != nil {
		t.Fatal(err)
	}

	s := NewSession(db)
	ps := NewPersonService(s)

	t.Log("Delete on these cases will return none affected")
	userids := []int{
		100000,
		200000,
	}
	for _, id := range userids {
		sid := strconv.Itoa(id)
		t.Run(sid, func(t *testing.T) {
			k := sky.PersonID(id)
			err := ps.Delete(k)
			if err != sky.ErrNoneAffected {
				t.Error(err)
			}
		})
	}
}
