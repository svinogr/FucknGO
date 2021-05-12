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
	setupShopHandlers(&f)
	testPanicHendler(&f)
	setupWebInterfaceHandler(&f)

	return f
}

// setupWebInterfaceHandler setup handlers for web
func setupWebInterfaceHandler(f *fabricHandlers) {
	hand := handler.MyHandler{"/login", loginPage, http.MethodGet, false, handler.TypeWeb}
	hand2 := handler.MyHandler{"/mainpage", mainPage, http.MethodGet, true, handler.TypeWeb}
	hand3 := handler.MyHandler{"/serverpage", serverPage, http.MethodGet, true, handler.TypeWeb}
	hand4 := handler.MyHandler{"/newuserpage", newuser, http.MethodGet, false, handler.TypeWeb}
	hand5 := handler.MyHandler{"/acountpage", accountPage, http.MethodGet, true, handler.TypeWeb}
	hand6 := handler.MyHandler{"/newshoppage", newShopPage, http.MethodGet, true, handler.TypeWeb}
	hand7 := handler.MyHandler{"/changeshoppage/{id}", changeShopPage, http.MethodGet, true, handler.TypeWeb}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3, &hand4, &hand5, &hand6, &hand7)
}

// setupUserHandlers setup handlers for actions with user
func setupUserHandlers(f *fabricHandlers) {
	// TODO установлено времено false потом переделать на плюс еще один столбик в таблице юзеры is active и активацией через email
	hand := handler.MyHandler{"/user", userApi, http.MethodPost, false, handler.TypeWeb}
	hand2 := handler.MyHandler{"/user/{id}", userApi, http.MethodDelete, true, handler.TypeWeb}

	f.Handlers = append(f.Handlers, &hand, &hand2)
}

// setupAuthHandlers setup handlers for actions with auth
func setupAuthHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/auth", auth, http.MethodPost, false, handler.TypeApi}
	hand4 := handler.MyHandler{"/logout", logOut, http.MethodGet, true, handler.TypeApi}
	hand3 := handler.MyHandler{"/auth/refresh-tokens", refreshToken, http.MethodPost, false, handler.TypeApi}

	f.Handlers = append(f.Handlers, &hand, &hand3, &hand4)
}

// setupShopHandlers setup handlers for actions with shopApi
func setupShopHandlers(f *fabricHandlers) {
	hand := handler.MyHandler{"/shop", shopApi, http.MethodPost, true, handler.TypeWeb}
	hand2 := handler.MyHandler{"/shop/{id}", shopApi, http.MethodDelete, true, handler.TypeWeb}
	hand3 := handler.MyHandler{"/shop/{id}", shopApi, http.MethodPut, true, handler.TypeWeb}

	f.Handlers = append(f.Handlers, &hand, &hand2, &hand3)
}

// setupServerHandlers setup handlers for actions with serverApi
func setupServerHandlers(f *fabricHandlers) {
	hand2 := handler.MyHandler{"/server", serverApi, http.MethodGet, true, handler.TypeWeb}
	hand3 := handler.MyHandler{"/server", serverApi, http.MethodPost, true, handler.TypeWeb}
	hand4 := handler.MyHandler{"/server/{id}", serverApi, http.MethodDelete, true, handler.TypeWeb}

	f.Handlers = append(f.Handlers, &hand2, &hand3, &hand4)
}

// test handler for imitation panic
func testPanicHendler(f *fabricHandlers) {
	hand := handler.MyHandler{"/panic", Panic, http.MethodGet, false, handler.TypeApi}
	f.Handlers = append(f.Handlers, &hand)
}
