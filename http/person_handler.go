package http

import (
	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

// we can have more data inside, funcs, db handlers, encoders, ...
type PersonHandler struct {
	*httprouter.Router
	PersonService sky.PersonService
}

func NewPersonHandler(s sky.PersonService) *PersonHandler {
	h := &PersonHandler{
		Router:        httprouter.New(),
		PersonService: s,
	}

	h.GET("/api/persons/:id", h.handleGetPerson)
	h.DELETE("/api/persons/:id", h.handleDeletePerson)
	h.PATCH("/api/persons/:id", h.handlePatchPerson)
	h.POST("/api/persons", h.handlePostPerson)

	return h
}
