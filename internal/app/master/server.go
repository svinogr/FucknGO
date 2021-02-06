package master

import (
	"FucknGO/internal/app/config"
	"fmt"
	"log"
	"net/http"
)

func RunServer() {
	conf := config.Get()

	// create new server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.MasterServer.Host, conf.MasterServer.Port),
		Handler: route(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while starting server, info: %v", err)
	}
}
