package server

import (
	"fmt"
	"net/http"
)

type Handler struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type FabricHandlers struct {
	Handlers []Handler
}

// NewFabric constructs new FabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() FabricHandlers {
	f := FabricHandlers{}
	f.Handlers = append(f.Handlers, Handler{"/", mainPage})

	return f
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}
