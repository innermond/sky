package http

import (
	"encoding/json"
	"net/http"

	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

func (h TokenHandler) handlePostToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		req    postTokenRequest
		tokstr string
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, sky.ErrInvalidJson, http.StatusBadRequest)
		return
	}
	k := sky.ApiKey(req.ApiKey)
	switch tokstr, err = h.TokenService.Create(k); err {
	case nil:
		encodeJson(w, &postTokenResponse{Token: tokstr})
	default:
		Error(w, err, http.StatusInternalServerError)
	}
}

type postTokenRequest struct {
	ApiKey string `json:"apikey"`
}

type postTokenResponse struct {
	Token string `json:"token"`
	errorResponse
}
