package server

import (
	"FucknGO/log"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
)

type fabricServers struct {
	servers []server
}

var instance *fabricServers
var once sync.Once

// FabricServer construct singleton
func FabricServer() (*fabricServers, error) {
	var err error
	once.Do(func() {
		instance = &fabricServers{}
		instance.servers = make([]server, 10)
	})

	return instance, err
}

// GetNewMasterServer creates and returns new master servers
func (f *fabricServers) GetNewMasterServer(address string, port string, staticResource string) *server {
	idServer := time.Now().Unix()

	server := server{}
	server.setup(address, port, staticResource, uint64(idServer))

	f.servers = append(f.servers, server)
	setupStaticResource(staticResource, server)

	fabric := NewFabric()
	for _, e := range fabric.Handlers {
		server.mux.HandleFunc(e.GetHandler().Path, e.GetHandler().HandlerFunc)
	}

	return &server
}

// GetNewSlaveServer creates and returns new slave  servers
func (f *fabricServers) GetNewSlaveServer(address string, port string, staticResource string) (*server, error) {

	for _, el := range f.servers {

		if el.port == port {
			return nil, errors.New("port is uses yet")
		}
	}

	idServer := time.Now().Unix()

	server := server{}
	server.setup(address, port, staticResource, uint64(idServer))

	f.servers = append(f.servers, server)
	setupStaticResource(staticResource, server)

	return &server, nil
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
