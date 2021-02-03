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
	hand2 := Handler.Handler{"/s", newServer}

	f.Handlers = append(f.Handlers, &hand, &hand2)

	return f
}

//main page
func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}

// newServer creates new server with input parameters
func newServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	port := query.Get("port")
	staticPath := query.Get("staticPath")

	if port != "" && staticPath != "" {
		server := Server{}
		go server.Start("127.0.0.1:"+port, staticPath)
		fmt.Fprint(w, "new server is run on port= "+port+"with static resource= "+staticPath)
	} else {
		fmt.Fprint(w, "invalid parameters")
	}
}
