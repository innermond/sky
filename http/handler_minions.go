package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/innermond/sky/sky"
)

// utilities:
func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{}"))
}

func Error(w http.ResponseWriter, err error, code int) {
	log.Printf("http error %s (code=%d)", err, code)
	if code == http.StatusInternalServerError {
		err = sky.ErrInternal
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

type errorResponse struct {
	Err string `json:"err,omitempty"`
}

func encodeJson(w http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError)
	}
}
