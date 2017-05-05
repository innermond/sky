package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

// we can have more data inside, funcs, db handlers, encoders, ...
type PersonHandler struct {
	*httprouter.Router
	PersonService sky.PersonService
}

func NewPersonHandler(s sky.PersonService) *PersonHandler {
	h := &PersonHandler{
		Router:        httprouter.New(),
		PersonService: s,
	}

	h.GET("/api/persons/:id", h.handleGetPerson)
	h.POST("/api/persons", h.handlePostPerson)

	return h
}

func (h PersonHandler) handleGetPerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pid, err := strconv.Atoi(ps.ByName("id"))
	log.Printf("get person %s", pid)
	p, err := h.PersonService.Get(sky.PersonID(pid))
	if err != nil {
		Error(w, err, http.StatusInternalServerError)
	} else if p == nil {
		NotFound(w)
	} else {
		encodeJson(w, &getPersonResponse{Person: *p})
	}
	// echo back url parameters
	s := fmt.Sprintf("%v %s", r.URL.Query(), pid)
	w.Write([]byte(s))
}

func (h PersonHandler) handlePostPerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var p postPersonRequest
}

type postPersonRequest struct {
	Person *sky.Person `json:"person,omitempty"`
}

type getPersonResponse struct {
	Person sky.Person `json:"person,omitempty"`

	errorResponse
}
