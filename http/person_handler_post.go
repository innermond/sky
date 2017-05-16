package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/innermond/sky/fail"
	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func (h PersonHandler) handlePostPerson(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		req postPersonRequest
		lid sky.PersonID
		err error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = errors.Wrap(err, "json decode")
		Error(w, err, http.StatusBadRequest)
		return
	}

	p := req.Person
	lid, err = h.PersonService.Create(p)
	// ok
	if err == nil {
		encodeJson(w, &postPersonResponse{Lid: lid})
		return
	}
	// errors
	switch err.(type) {
	case fail.Mistakes:
		Error(w, err.(fail.Mistakes), http.StatusPreconditionFailed)
		return
	}

	Error(w, err, http.StatusInternalServerError)
}

type postPersonRequest struct {
	Person sky.Person `json:"person,omitempty"`
}

type postPerson postPersonRequest

func (pq *postPersonRequest) UnmarshalJSON(data []byte) error {
	var preq = postPerson{}
	err := json.Unmarshal(data, &preq)
	if err == nil {
		preq.Person.Longname = strings.TrimSpace(preq.Person.Longname)
		*pq = postPersonRequest(preq)
	}
	return err
}

type postPersonResponse struct {
	Lid sky.PersonID `json:"lid,omitempty"`
	errorResponse
}
