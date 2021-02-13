package main

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/log"
	"FucknGO/internal/app/master"
)

func main() {
	conf := config.Get()

	if !conf.Basic.Debug {
		log.SetOutputFile()
	}
	logger := log.GetLogger()

	logger.Printf("Basic: %+v\n", conf.Basic)
	logger.Printf("MS: %+v\n", conf.MasterServer)

	go master.RunServer()
	logger.Printf("Successfully started master server on http://%s:%d\n",
		conf.MasterServer.Host, conf.MasterServer.Port)

	select {}
}
