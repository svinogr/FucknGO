package server

import (
	"FucknGO/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	Config config.Config
}

// Start starts server with settings
func (s *Server) Start() {
	// TODO вставить обработку ошибку отсутсвия конфига

	s.SetupHttpHandlers()

	s.setupStaticResource()

	go s.runServer()

	select {}
}

func (s *Server) runServer() {
	port := s.Config.JsonStr.ServerConfig.Port
	err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(port)), nil)

	if err != nil {
		log.Fatal(err)
	}
}

//SetupHttpHandlers set handlers for http.HandleFunc
func (s *Server) SetupHttpHandlers() {
	fabric := NewFabric()

	for i, e := range fabric.Handlers {
		fmt.Println(i)
		http.HandleFunc(e.GetHandler().Path, e.GetHandler().HandlerFunc)
	}
}

// setupStaticResource creates a directory named path: "./app/server/ui/static",
// along with any necessary parents, and returns nil,
// or else returns an error.
// If path is already a directory, setupStaticResource does nothing
// and returns nil.
func (s *Server) setupStaticResource() {

	path := s.Config.JsonStr.UiConfig.WWW.Static

	_, err := os.Stat(path)

	if err != nil {
		err = os.MkdirAll(path, 0777)

		if err != nil {
			log.Fatal(err)
		}
	}

	fileServer := http.FileServer(http.Dir(path))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}
