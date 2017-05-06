package http

import (
	"encoding/json"
	"log"
	"net/http"

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

	switch err = h.PersonService.Modify(p); err {
	case nil:
		encodeJson(w, &patchPersonResponse{})
	case sky.ErrPersonValid:
		log.Println(err)
		Error(w, sky.ErrPersonValid, http.StatusBadRequest)
	default:
		Error(w, err, http.StatusInternalServerError)
	}
}

type patchPersonRequest struct {
	Person sky.Person `json:"person,omitempty"`
}

type patchPersonResponse struct {
	errorResponse
}
