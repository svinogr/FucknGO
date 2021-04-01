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
	setupWebInterfaceHandler(&f)
	return f
}

// setupWebInterfaceHandler setup handlers for web
func setupWebInterfaceHandler(f *fabricHandlers) {
	hand := handler.MyHandler{
		Path:          "/login",
		HandlerFunc:   loginPage,
		Method:        http.MethodGet,
		NeedAuthToken: false,
	}
	hand2 := handler.MyHandler{
		Path:          "/mainpage",
		HandlerFunc:   mainPage,
		Method:        http.MethodGet,
		NeedAuthToken: false,
	}
	hand3 := handler.MyHandler{
		Path:          "/logout",
		HandlerFunc:   logOut,
		Method:        http.MethodGet,
		NeedAuthToken: false,
	}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3)
}

// setupUserHandlers setup handlers for actions with user
func setupUserHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/user", user, http.MethodPost, false}

	f.Handlers = append(f.Handlers, &hand)
}

// setupAuthHandlers setup handlers for actions with auth
func setupAuthHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/auth", auth, http.MethodPost, false}
	hand2 := handler.MyHandler{"/log", logPage, http.MethodGet, false}
	hand3 := handler.MyHandler{"/auth/refreshtoken", refreshToken, http.MethodPost, false}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3)
}

// setupServerHandlers setup handlers for actions with server
func setupServerHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/", MainPage, http.MethodGet, false}
	hand2 := handler.MyHandler{"/server", Server, http.MethodGet, true}
	hand3 := handler.MyHandler{"/server", Server, http.MethodPost, true}
	hand4 := handler.MyHandler{"/server/{id}", Server, http.MethodDelete, true}
	hand5 := handler.MyHandler{"/connect", Connect, http.MethodGet, false}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5)
}

// test handler for imitation panic
func testPanicHendler(f *fabricHandlers) {
	hand := handler.MyHandler{"/panic", Panic, http.MethodGet, false}
	f.Handlers = append(f.Handlers, &hand)
}
