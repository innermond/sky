package http

import (
	"fmt"
	"net/http"
	"strings"
)

type AllServicesHandler struct {
	PersonHandler *PersonHandler
}

func (h *AllServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.Path // /api/entity... we are interested in entity
	parts := strings.Split(uri, "/")
	if len(parts) < 3 {
		Error(w, fmt.Errorf("%v", "url malformed"), http.StatusBadRequest)
		return
	}

	resource := parts[2]

	// check presence auth token for entire api's endpoints excepts "authenticate"
	tokenName := "Autorization"
	if resource != "authenticate" && "" == r.Header.Get(tokenName) {
		NotAuthenticated(w)
		return
	}

	switch resource {
	case "persons":
		h.PersonHandler.ServeHTTP(w, r)
	default:
		NotFound(w)
		return
	}
}
