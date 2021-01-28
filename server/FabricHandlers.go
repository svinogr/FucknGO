package server

import (
	"FucknGO/server/Handler"
	"fmt"
	"net/http"
)

type FabricHandlers struct {
	Handlers []Handler.HandlerInterface
}

// NewFabric constructs new FabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() FabricHandlers {
	f := FabricHandlers{}
	hand := Handler.Handler{"/", mainPage}
	f.Handlers = append(f.Handlers, &hand)

	return f
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}
