package server

import (
	"FucknGO/log"
	"net/http"
	"os"
	"sync"
)

type fabricServers struct {
	server []server
}

var instance *fabricServers
var once sync.Once

// FabricServer construct singleton
func FabricServer() (*fabricServers, error) {
	var err error
	once.Do(func() {
		instance = &fabricServers{}
		instance.server = make([]server, 10)
	})

	return instance, err
}

// GetNewMasterServer creates and returns new master server
func (f *fabricServers) GetNewMasterServer(address string, staticResource string) *server {
	server := server{}
	server.setup(address, staticResource)

	f.server = append(f.server, server)
	setupStaticResource(staticResource, server)

	fabric := NewFabric()
	for _, e := range fabric.Handlers {
		server.mux.HandleFunc(e.GetHandler().Path, e.GetHandler().HandlerFunc)
	}

	return &server
}

// GetNewSlaveServer creates and returns new slave  server
func (f *fabricServers) GetNewSlaveServer(address string, staticResource string) *server {
	server := server{}
	server.setup(address, staticResource)

	f.server = append(f.server, server)
	setupStaticResource(staticResource, server)

	return &server
}

//setupStaticResource set static dir for server
func setupStaticResource(staticResource string, server server) {
	_, err := os.Stat(staticResource)

	if err != nil {
		err = os.MkdirAll(staticResource, 0777)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}

	fileServer := http.FileServer(http.Dir(staticResource))

	server.mux.Handle("/static/", http.StripPrefix("/static", fileServer))
}
