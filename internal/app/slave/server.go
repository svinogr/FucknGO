package slave

import (
	"FucknGO/internal/app/config"
	"FucknGO/pkg/requests"
	"fmt"
	"log"
	"net/http"
)

var Slaves []*http.Server

func RunServer(params *requests.SlaveParams) {
	conf := config.Get()

	// create new server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.MasterServer.Host, params.Port),
		Handler: route(params.StaticDir),
	}

	Slaves = append(Slaves, &server)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting server, info: %v", err)
	}
}

func GetSlaves() []*http.Server {
	return Slaves
}
