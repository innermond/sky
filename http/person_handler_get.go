package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

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

type getPersonResponse struct {
	Person sky.Person `json:"person,omitempty"`
	errorResponse
}
