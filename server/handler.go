package server

import (
	"net/http"

	"github.com/elliottpolk/peppermint-sparkles/internal/respond"
)

const (
	path string = "/api/v4/secrets"

	appParam  string = "app_name"
	envParam  string = "env"
	userParam string = "username"
	idParam   string = "uuid"
)

type Handler struct {
}

func Handle(mux *http.ServeMux, h *Handler) *http.ServeMux {
	// TODO:
	return mux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO:
	switch r.Method {
	default:
		respond.MethodNotAllowed(w)
	}
}
