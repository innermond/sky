package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/innermond/sky/sky"
)

type AllServicesHandler struct {
	Auth          sky.Authenticator
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
	tokenstr := r.Header.Get(tokenName)
	if resource != "authenticate" && "" == tokenstr {
		NotAuthenticated(w)
		return
	}
	// try to parse a jwt token
	tokenstr = strings.TrimSpace(tokenstr)
	if !strings.HasPrefix(tokenstr, "Bearer") {
		NotAuthenticated(w)
		return
	}
	tk := tokenstr[7:]
	err := h.Auth.Authenticate(tk)
	if err != nil {
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
