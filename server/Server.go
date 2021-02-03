package server

import (
	"FucknGO/log"
	"net/http"
	"os"
)

type Server struct {
}

// Start creates and starts server with settings
func (s *Server) Start(address string, staticResource string) {
	mux := http.ServeMux{}

	s.setupStaticResource(staticResource, &mux)
	s.setupHttpHandlers(&mux)
	s.runServer(address, &mux)
}

// runServer run server on address by value address
func (s *Server) runServer(address string, mux *http.ServeMux) {
	server := http.Server{Addr: address, Handler: mux}

	err := server.ListenAndServe()

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

//setupHttpHandlers set handlers for http.HandleFunc
func (s *Server) setupHttpHandlers(mux *http.ServeMux) *http.ServeMux {
	fabric := NewFabric()

	for _, e := range fabric.Handlers {
		mux.HandleFunc(e.GetHandler().Path, e.GetHandler().HandlerFunc)
	}
	return mux
}

// setupStaticResource creates a directory named by value staticResource ,
// along with any necessary parents,
func (s *Server) setupStaticResource(staticResource string, mux *http.ServeMux) *http.ServeMux {
	path := staticResource

	_, err := os.Stat(path)

	if err != nil {
		err = os.MkdirAll(path, 0777)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}

	fileServer := http.FileServer(http.Dir(staticResource))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
