package main

import (
	"FucknGO/config"
	"FucknGO/server"
	"log"
)

func main() {
	conf := config.Config{}
	_, err := conf.ReadConfig(config.Path)

	if err != nil {
		log.Fatal(err)
	}

	server := server.Server{Config: conf}

	server.Start()

	//fmt.Print(config.JsonStr.ServerConfig.JwtLifeTimeDays)

	//здесь запускаем сервер
}
