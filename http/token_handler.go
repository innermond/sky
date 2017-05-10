package http

import (
	"github.com/innermond/sky/sky"
	"github.com/julienschmidt/httprouter"
)

type TokenHandler struct {
	*httprouter.Router
	TokenService sky.TokenService
}

func NewTokenHandler(s sky.TokenService) *TokenHandler {
	h := &TokenHandler{
		Router:       httprouter.New(),
		TokenService: s,
	}

	h.POST("/api/tokens", h.handlePostToken)

	return h
}
