package slave

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/log"
	"FucknGO/pkg/requests"
	"fmt"
	"net/http"
	"time"
)

var Slaves map[int64]*http.Server

func RunServer(params *requests.SlaveParams) {
	conf := config.Get()
	logger := log.GetLogger()

	// create new server
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.MasterServer.Host, params.Port),
		Handler: route(params.StaticDir),
	}

	slaves := GetSlaves()
	slaves[time.Now().Unix()] = &server

	err := server.ListenAndServe()
	if err != nil {
		logger.Printf("Error while starting server, info: %v", err)
	}
}

func GetSlaves() map[int64]*http.Server {
	if Slaves == nil {
		Slaves = make(map[int64]*http.Server)
	}
	return Slaves
}
