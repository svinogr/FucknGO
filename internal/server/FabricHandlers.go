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
	setupUserHandlers(&f)
	testPanicHendler(&f)

	return f
}

func setupUserHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/user", user, http.MethodPost, false}

	f.Handlers = append(f.Handlers, &hand)
}

func setupAuthHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/auth", auth, http.MethodPost, false}
	hand2 := handler.MyHandler{"/log", logPage, http.MethodGet, false}

	f.Handlers = append(f.Handlers, &hand, &hand2)
}

func setupServerHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/", MainPage, http.MethodGet, false}
	hand2 := handler.MyHandler{"/server", Server, http.MethodGet, true}
	hand3 := handler.MyHandler{"/server", Server, http.MethodPost, true}
	hand4 := handler.MyHandler{"/server/{id}", Server, http.MethodDelete, true}
	hand5 := handler.MyHandler{"/connect", Connect, http.MethodGet, false}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5)
}

func testPanicHendler(f *fabricHandlers) {
	hand := handler.MyHandler{"/panic", Panic, http.MethodGet, false}
	f.Handlers = append(f.Handlers, &hand)
}
