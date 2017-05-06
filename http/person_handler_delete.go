package http

import (
	"net/http"
	"strconv"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

func (h PersonHandler) handleDeletePerson(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_pid, err := strconv.Atoi(ps.ByName("id"))
	pid := sky.PersonID(_pid)
	err = h.PersonService.Delete(pid)
	switch err {
	case nil:
		encodeJson(w, &deletePersonResponse{Pid: pid})
	case sky.ErrNoneAffected:
		Error(w, err, http.StatusUnprocessableEntity)
	default:
		Error(w, err, http.StatusInternalServerError)
	}
}

type deletePersonResponse struct {
	Pid sky.PersonID `json:"pid,omitempty"`
	errorResponse
}
