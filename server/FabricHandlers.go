package server

import (
	"fmt"
	"net/http"
)

type FabricHandlers struct {
	Handlers []HandlerInterface
}

// NewFabric constructs new FabricHandlers and inflate handlers for http.HandleFunc
func NewFabric() FabricHandlers {
	f := FabricHandlers{}
	hand := Handler{"/", mainPage}
	f.Handlers = append(f.Handlers, &hand)

	return f
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}
