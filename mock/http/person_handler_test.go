package http

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/innermond/sky/auth"
	"github.com/innermond/sky/config"
	myhttp "github.com/innermond/sky/http"
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
	db := config.DB()
	verify := config.PublicKey()
	if err := config.Err(); err != nil {
		panic(err)
	}
	a := auth.NewAuthenticator(verify)
	// session
	s := mysql.NewSession(db)
	// services
	personService := mysql.NewPersonService(s)
	all := &myhttp.AllServicesHandler{
		AllRoutesAuth: a,
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
			req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer([]byte(uc.data)))
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE0OTQ2NzUxNDEsIklkIjoxfQ.1k7QTZlpBAIQ6NqQlH3J7vz6CeuTPiV-E3crW1rExyrvtpW8alELvE1gSNAACaHY22y0_2GJVJmydEJI3zP8LOpvHj7aH_utQw2hGajzDtdK2BnlNwoKAmEUEzi70-eHd2AzDsvj49AVBHQbN_vh8Hd49Q5grM35iYHXRIh82VY"))
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
