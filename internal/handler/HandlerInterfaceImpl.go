package handler

import "net/http"

type MyHandler struct {
	Path        string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
	Method      string
}

func (h *MyHandler) GetHandler() *MyHandler {
	return h
}
