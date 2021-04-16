package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	id             uint64
	mux            mux.Router
	address        string
	port           string
	staticResource string
	server         http.Server
	isSlave        bool
}

func (s *Server) Port() string {
	return s.port
}

func (s *Server) StaticResource() string {
	return s.staticResource
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) Id() uint64 {
	return s.id
}

// Setup creates and starts serverApi with settings
func (s *Server) setup(address string, port string, id uint64, isSlave bool) {
	s.isSlave = isSlave
	s.mux = *mux.NewRouter()
	s.address = address
	s.port = port
	s.id = id

	setupStaticResource(s)
	setupHandlers(s)
}
