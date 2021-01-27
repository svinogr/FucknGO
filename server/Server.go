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
func (s *Server) Start() error {
	// TODO вставить обработку ошибку отсутсвия конфига
	port := s.Config.JsonStr.ServerConfig.Port
	s.SetupHttpHandlers()

	err := s.setupStaticResource()

	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(strconv.Itoa(int(port)), nil)

	if err != nil {
		return err
	} else {
		return nil
	}
}

//SetupHttpHandlers set handlers for http.HandleFunc
func (s *Server) SetupHttpHandlers() {
	handlers := FabricHandlers{}
	fabric := handlers.NewFabric()

	for i, e := range fabric.Handlers {
		fmt.Println(i)
		http.HandleFunc(e.Path, e.Handler)
	}
}

// setupStaticResource creates a directory named path: "./app/server/ui/static",
// along with any necessary parents, and returns nil,
// or else returns an error.
// If path is already a directory, setupStaticResource does nothing
// and returns nil.
func (s *Server) setupStaticResource() error {

	path := s.Config.JsonStr.UiConfig.WWW.Static

	_, err := os.Stat(path)

	if err != nil {
		err = os.MkdirAll(path, 0777)

		if err != nil {
			return err
		}
	}

	fileServer := http.FileServer(http.Dir(path))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	return nil
}

// startServer starts server on port 8080
func startServer() {
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error while starting serverErrors: %v", err)
	}
}
