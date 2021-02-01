package server

import (
	"FucknGO/config"
	"FucknGO/log"
	"net/http"
	"os"
)

type Server struct {
	Config config.Config
}

func init() {
	SetupHttpHandlers()
}

// Start starts server with settings
func (s *Server) Start(address string, staticResource string) {
	s.setupStaticResource(staticResource)

	s.runServer(address)
}

func (s *Server) runServer(address string) {
	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

//SetupHttpHandlers set handlers for http.HandleFunc
func SetupHttpHandlers() {
	fabric := NewFabric()

	for _, e := range fabric.Handlers {
		http.HandleFunc(e.GetHandler().Path, e.GetHandler().HandlerFunc)
	}
}

// setupStaticResource creates a directory named path: "./app/server/ui/static",
// along with any necessary parents, and returns nil,
// or else returns an error.
// If path is already a directory, setupStaticResource does nothing
// and returns nil.
func (s *Server) setupStaticResource(staticResource string) {

	path := staticResource

	_, err := os.Stat(path)

	if err != nil {
		err = os.MkdirAll(path, 0777)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}

	fileServer := http.FileServer(http.Dir(staticResource))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}
