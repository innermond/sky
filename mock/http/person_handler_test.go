package http

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/innermond/sky/config"
	myhttp "github.com/innermond/sky/http"
	"github.com/innermond/sky/jwt"
	"github.com/innermond/sky/mysql"
)

func init() {
	noLog := flag.Bool("nolog", true, "disable log output")
	flag.Parse()
	if !*noLog {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.LstdFlags)
}

// act like main() basically it is an adjusted copy paste of main.main()
func minor() *httptest.Server {
	dns := "root:M0b1d1c3@tcp(localhost:3306)/printoo"
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	// authenticator
	verify, err := config.PublicKey()
	if err != nil {
		panic(err)
	}
	a := jwt.NewAuthenticator(verify)
	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	all := &myhttp.AllServicesHandler{
		Auth:          a,
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
		data     string
		expected int
	}{
		// longname is not valid
		{`{"person":{"id":1,"longname":"ga"}}`, 412},
		{`{"person":{"id":1,"longname":"ga is more than 10 characters length"}}`, 412},
		{`{"person":{"id":1,"longname":"I try to\nmake it fun"}}`, 412},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons"
	for _, uc := range uses {
		t.Run(uc.data, func(t *testing.T) {
			req, err := http.NewRequest("PATCH", urlStr, bytes.NewBuffer([]byte(uc.data)))
			fatalif(err, t)
			res, err := http.DefaultClient.Do(req)
			fatalif(err, t)
			if res.StatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
	}
}

func TestDeleteSinglePerson(t *testing.T) {
	uses := []struct {
		id       string
		expected int
	}{
		// deleting none returns code 422 UnprocessableEntity
		{"100", 422},
		{"aaa", 422},
		{"0", 422},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons/"
	for _, uc := range uses {
		t.Run(uc.id, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", urlStr+uc.id, nil)
			fatalif(err, t)
			res, err := http.DefaultClient.Do(req)
			fatalif(err, t)
			if res.StatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
	}
}

func TestPatchSinglePerson(t *testing.T) {
	uses := []struct {
		data     string
		expected int
	}{
		// longname is not valid
		{`{"person":{"id":1,"longname":"ga"}}`, 412},
		{`{"person":{"id":1,"longname":"ga is more than 10 characters length"}}`, 412},
		{`{"person":{"id":1,"longname":"I try to\nmake it fun"}}`, 412},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons"
	for _, uc := range uses {
		t.Run(uc.data, func(t *testing.T) {
			req, err := http.NewRequest("PATCH", urlStr, bytes.NewBuffer([]byte(uc.data)))
			fatalif(err, t)
			res, err := http.DefaultClient.Do(req)
			fatalif(err, t)
			if res.StatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
	}
}

func TestPerson_checkToken(t *testing.T) {
	uses := []struct {
		id       string
		expected int
	}{
		// forbidden http code
		{"1", 403},
		{"99", 403},
		{"100", 403},
	}
	srv := minor()
	defer srv.Close()
	urlStr := srv.URL + "/api/persons/"
	for _, uc := range uses {
		t.Run(uc.id, func(t *testing.T) {

			req, err := http.NewRequest("GET", urlStr+uc.id, nil)
			fatalif(err, t)
			tk := "fake.token.baby"
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %w", tk))
			res, err := http.DefaultClient.Do(req)
			fatalif(err, t)
			if res.StatusCode != uc.expected {
				t.Errorf("status code expected %d got %d", uc.expected, res.StatusCode)
			}
		})
		t.Run("authenticate", func(t *testing.T) {
			res, err := http.Get(srv.URL + "/authenticate")
			fatalif(err, t)
			if res.StatusCode == 403 {
				t.Errorf("status code unexpected %d", res.StatusCode)
			}

		})
	}
}
