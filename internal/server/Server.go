package server

import (
	"FucknGO/log"
	"net/http"
)

type server struct {
	id             uint64
	mux            http.ServeMux
	address        string
	port           string
	staticResource string
	server         http.Server
}

func (s *server) Port() string {
	return s.port
}

func (s *server) StaticResource() string {
	return s.staticResource
}

func (s *server) Address() string {
	return s.address
}

func (s *server) Id() uint64 {
	return s.id
}

// Setup creates and starts server with settings
func (s *server) setup(address string, port string, staticResource string, id uint64) {
	s.mux = http.ServeMux{}
	s.address = address
	s.port = port
	s.staticResource = staticResource
	s.id = id
}

// runServer run servers
func (s *server) RunServer() {
	s.server = http.Server{Addr: s.address + ":" + s.port, Handler: &s.mux}

	err := s.server.ListenAndServe()

	if err != nil {
		if err.Error() != "http: Server closed" {
			log.NewLog().Fatal(err)
		}

		log.NewLog().PrintError(err)
	}
}
