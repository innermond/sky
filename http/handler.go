package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type AllServicesHandler struct {
	PersonHandler *PersonHandler
}

func (h *AllServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path // /api/entity... we are interested in entity
	parts := strings.Split(uri, "/")
	log.Println(uri, len(parts))
	if len(parts) < 3 {
		Error(w, fmt.Errorf("%v", "url malformed"), http.StatusBadRequest)
		return
	}

	switch parts[2] {
	case "persons":
		h.PersonHandler.ServeHTTP(w, r)
	default:
		NotFound(w)
		return
	}
}
