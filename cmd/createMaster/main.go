package main

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/db"
	"FucknGO/internal/app/master"
	"log"
)

func main() {
	conf := config.Get()
	log.Printf("DB: %+v\n", conf.DB)
	log.Printf("MS: %+v\n", conf.MasterServer)

	orm := db.GetDB()
	log.Printf("Successfully connect to DB %v\n", orm)

	go master.RunServer()
	log.Printf("Successfully started master server on http://%s:%d\n",
		conf.MasterServer.Host, conf.MasterServer.Port)

	select {}
}
