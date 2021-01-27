package main

import (
	"FucknGO/config"
	"fmt"
	"log"
)

func main() {
	conf := config.Config{}
	config, err := conf.ReadConfig(config.Path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(config.JsonStr.ServerConfig.JwtLifeTimeDays)

	//здесь запускаем сервер
}
