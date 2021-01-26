package main

import (
	"FucknGO/config"
	"fmt"
	"log"
)

func main() {
	conf := config.Config{}
	config, err := conf.ReadConfig("./config/config.json")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(config.ServerConfig.Port)
	fmt.Print(config.ServerConfig.MiddlewareTimeout)

	//здесь запускаем сервер
}
