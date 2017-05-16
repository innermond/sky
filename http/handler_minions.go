package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/innermond/sky/fail"
	"github.com/innermond/sky/sky"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

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
	logput, has := trace(err)
	if has {
		log.Printf("trace %s (code=%d)", logput, code)
	} else {
		log.Printf("http error %s (code=%d)", err, code)
	}
	//err = errors.Cause(err)
	switch code {
	case http.StatusInternalServerError:
		err = sky.ErrInternal
	case http.StatusPreconditionFailed:
		if merrs, ok := err.(fail.Mistakes); ok {
			err = merrs
			break
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

func trace(err error) (out string, exists bool) {
	var terr stackTracer
	if terr, exists = err.(stackTracer); exists {
		for _, f := range terr.StackTrace() {
			out += "\n" + fmt.Sprintf("%+s:%d", f)
		}
	}
	return
}
