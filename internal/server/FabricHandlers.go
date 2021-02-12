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
	hand2 := handler.Handler{"/server", Server}
	hand3 := handler.Handler{"/connect", Connect}
	hand4 := handler.Handler{"/server/all", GetAllServers}
	hand5 := handler.Handler{"/serverkill", StopServerById}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5)

	return f
}
