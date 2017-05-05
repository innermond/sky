package http

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	myhttp "github.com/innermond/sky/http"
	"github.com/innermond/sky/mysql"
)

// act like main() basically it is an adjusted copy paste of main.main()
func minor() *httptest.Server {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags)

	dns := "root:M0b1d1c3@tcp(localhost:3306)/printoo"
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	all := &myhttp.AllServicesHandler{
		PersonHandler: myhttp.NewPersonHandler(personService),
	}
	srv := httptest.NewServer(all)
	return srv
}

func fatalif(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSinglePerson(t *testing.T) {
	uses := []struct {
		id       string
		expected int
	}{
		// nonexistent
		{"99", 404},
		{"100", 404},
		// malformed
		//{"a", 422},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons/"
	for _, uc := range uses {
		t.Run(uc.id, func(t *testing.T) {
			res, err := http.Get(urlStr + uc.id)
			fatalif(err, t)
			if res.StatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
	}
}

func TestPostSinglePerson(t *testing.T) {
	uses := []struct {
		longname string
	}{
		// nonexistent
		{"gabi1"},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons"
	for _, uc := range uses {
		t.Run(uc.longname, func(t *testing.T) {
			res, err := http.Get(urlStr + uc.id)
			fatalif(err, t)
			if resStatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
	}
}
