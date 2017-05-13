package http

import (
	"encoding/json"
	"net/http"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

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

/*func (pq *postPersonRequest) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, pq)
	if err == nil {
		pq.Person.Longname = "am modified"
	}
	return err
}*/

type postPersonResponse struct {
	Lid sky.PersonID `json:"lid,omitempty"`
	errorResponse
}
