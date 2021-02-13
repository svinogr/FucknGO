package master

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/log"
	"fmt"
	"net/http"
)

func RunServer() {
	conf := config.Get()
	logger := log.GetLogger()

	// create new server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.MasterServer.Host, conf.MasterServer.Port),
		Handler: route(),
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Printf("Error while starting server, info: %v", err)
	}
}
