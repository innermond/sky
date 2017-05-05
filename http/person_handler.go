package http

import (
	"encoding/json"
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
}

func (h PersonHandler) handlePostPerson(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		req postPersonRequest
		lid sky.PersonID
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, sky.ErrInvalidJson, http.StatusBadRequest)
		return
	}

	p := req.Person

	switch lid, err = h.PersonService.Create(p); err {
	case nil:
		encodeJson(w, &postPersonResponse{Lid: lid})
	case sky.ErrPersonValid:
		Error(w, sky.ErrPersonValid, http.StatusBadRequest)
	default:
		Error(w, err, http.StatusInternalServerError)
	}
}

type postPersonRequest struct {
	Person sky.Person `json:"person,omitempty"`
}

type postPersonResponse struct {
	Lid sky.PersonID `json:"lid,omitempty"`
	errorResponse
}

type getPersonResponse struct {
	Person sky.Person `json:"person,omitempty"`
	errorResponse
}
