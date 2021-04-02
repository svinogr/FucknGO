package handler

import "net/http"

type TypeRequest int

const (
	TypeWeb = iota
	TypeApi
)

type MyHandler struct {
	Path          string
	HandlerFunc   func(w http.ResponseWriter, r *http.Request)
	Method        string
	NeedAuthToken bool
	TypeRequest   TypeRequest
}

func (h *MyHandler) GetHandler() *MyHandler {
	return h
}
