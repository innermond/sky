package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/innermond/sky/fail"
	"github.com/innermond/sky/sky"
)

// utilities:
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{}"))
}
func NotAuthenticated(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("{}"))
}

func Error(w http.ResponseWriter, err error, code int) {
	log.Printf("http error %s (code=%d)", err, code)
Code:
	switch code {
	case http.StatusInternalServerError:
		err = sky.ErrInternal
	case http.StatusPreconditionFailed:
		if merrs, ok := err.(fail.Mistakes); ok {
			err = merrs
			break Code
		}
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err})
}

type errorResponse struct {
	Err interface{} `json:"err,omitempty"`
}

func encodeJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}
