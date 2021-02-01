package server

import (
	"FucknGO/config"
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
	//	f.Handlers = append(f.Handlers, &hand2)

	return f
}

//main page
func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}

func newServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	port := query.Get("port")

	config, _ := config.GetConfig()
	server := Server{Config: *config}
	go server.Start("127.0.0.1:"+port, "./ui/web/slave/")

	/*	query := r.URL.Query()
		port := query.Get("port")

		re, err := http.Get("localhost:" + port)

		fmt.Fprint(w, "port is busy")
		fmt.Print(err)
		fmt.Print(re)
		return

		go http.ListenAndServe(":"+port, nil)
		fmt.Fprint(w, query)*/

}
