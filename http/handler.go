package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

type AllServicesHandler struct {
	PersonHandler *PersonHandler
}

func (h *AllServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path // /api/entity... we are interested in entity
	parts := strings.Split(uri, "/")
	log.Println(uri, len(parts))
	if len(parts) < 3 {
		Error(w, fmt.Errorf("%v", "url malformed"), http.StatusBadRequest)
		return
	}

	switch parts[2] {
	case "persons":
		h.PersonHandler.ServeHTTP(w, r)
	default:
		NotFound(w)
		return
	}
}

// we can have more data inside, funcs, db handlers, encoders, ...
type PersonHandler struct {
	*httprouter.Router
}

func NewPersonHandler() *PersonHandler {
	h := &PersonHandler{
		Router: httprouter.New(),
	}

	h.GET("/api/persons/:id", h.handleGetPerson)

	return h
}

func (h PersonHandler) handleGetPerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid := ps.ByName("id")
	log.Printf("get person %s", pid)
	// echo back url parameters
	s := fmt.Sprintf("%v %s", r.URL.Query(), pid)
	w.Write([]byte(s))
}

type getPersonResponse struct {
	Person sky.Person `json:"person,omitempty"`

	errorResponse
}