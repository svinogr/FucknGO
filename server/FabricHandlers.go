package server

import (
	"FucknGO/server/Handler"
	"fmt"
	"net/http"
)

type fabricHandlers struct {
	Handlers []Handler.HandlerInterface
}

// NewFabric constructs new fabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() fabricHandlers {
	f := fabricHandlers{}
	hand := Handler.Handler{"/", mainPage}
	f.Handlers = append(f.Handlers, &hand)

	return f
}

//main page
func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}
