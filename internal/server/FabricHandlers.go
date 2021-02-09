package server

import (
	"FucknGO/internal/handler"
)

type fabricHandlers struct {
	Handlers []handler.HandlerInterface
}

// NewFabric constructs new fabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() fabricHandlers {
	f := fabricHandlers{}

	hand := handler.Handler{"/", MainPage}
	hand2 := handler.Handler{"/s", NewServer}
	hand3 := handler.Handler{"/connect", Connect}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3)

	return f
}
