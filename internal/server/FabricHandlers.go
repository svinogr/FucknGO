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

	hand := handler.Handler{"/", MainPage, "GET"}
	hand2 := handler.Handler{"/server", Server, "GET"}
	hand3 := handler.Handler{"/server", Server, "POST"}
	hand4 := handler.Handler{"/server/{id}", Server, "DELETE"}
	hand5 := handler.Handler{"/connect", Connect, "GET"}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5)

	return f
}
