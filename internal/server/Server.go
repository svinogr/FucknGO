package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type server struct {
	id             uint64
	mux            mux.Router
	address        string
	port           string
	staticResource string
	server         http.Server
	isSlave        bool
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

// Setup creates and starts serverApi with settings
func (s *server) setup(address string, port string, id uint64, isSlave bool) {
	s.isSlave = isSlave
	s.mux = *mux.NewRouter()
	s.address = address
	s.port = port
	s.id = id

	setupStaticResource(s)
	setupHandlers(s)
}

// runServer run servers
func (s *server) RunServer() error {
	s.server = http.Server{Addr: s.address + ":" + s.port, Handler: handlers.LoggingHandler(os.Stdout, &s.mux)} //TODO настроить запись в файл

	return s.server.ListenAndServe()
}
