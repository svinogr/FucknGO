package server

import (
	"FucknGO/internal/handler"
	"net/http"
)

type fabricHandlers struct {
	Handlers []handler.HandlerInterface
}

// NewFabric constructs new fabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() fabricHandlers {
	f := fabricHandlers{}
	setupServerHandlers(&f)
	setupAuthHandlers(&f)

	return f
}

func setupAuthHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/auth", auth, http.MethodPost}
	hand2 := handler.MyHandler{"/log", logPage, http.MethodGet}

	f.Handlers = append(f.Handlers, &hand, &hand2)
}

func setupServerHandlers(f *fabricHandlers) {

	hand := handler.MyHandler{"/", MainPage, http.MethodGet}
	hand2 := handler.MyHandler{"/server", Server, http.MethodGet}
	hand3 := handler.MyHandler{"/server", Server, http.MethodPost}
	hand4 := handler.MyHandler{"/server/{id}", Server, http.MethodDelete}
	hand5 := handler.MyHandler{"/connect", Connect, http.MethodGet}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5)
}
