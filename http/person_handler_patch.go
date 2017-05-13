package http

import (
	"encoding/json"
	"net/http"

	"github.com/innermond/sky/fail"
	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

func (h PersonHandler) handlePatchPerson(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		req patchPersonRequest
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, sky.ErrInvalidJson, http.StatusBadRequest)
		return
	}

	p := req.Person

	err = h.PersonService.Modify(p)
	// Validation
	if err == nil {
		encodeJson(w, &patchPersonResponse{})
		return
	}
	if verr, ok := err.(*fail.Mistake); ok {
		Error(w, verr, http.StatusPreconditionFailed)
		return
	}
	Error(w, err, http.StatusInternalServerError)
}

type patchPersonRequest struct {
	Person sky.Person `json:"person,omitempty"`
}

type patchPersonResponse struct {
	errorResponse
}
