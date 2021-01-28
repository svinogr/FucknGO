package server

import "net/http"

type Handler struct {
	Path        string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func (h *Handler) GetHandler() *Handler {
	return h
}
