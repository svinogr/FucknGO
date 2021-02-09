package server

import (
	"FucknGO/log"
	"net/http"
)

type server struct {
	mux            *http.ServeMux
	address        string
	staticResource string
}

// Setup creates and starts server with settings
func (s *server) setup(address string, staticResource string) {
	s.mux = &http.ServeMux{}
	s.address = address
	s.staticResource = staticResource
}

// runServer run server
func (s *server) RunServer() {
	server := http.Server{Addr: s.address, Handler: s.mux}

	err := server.ListenAndServe()

	if err != nil {
		log.NewLog().Fatal(err)
	}
}
