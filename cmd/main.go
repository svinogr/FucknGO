package main

import (
	"FucknGO/config"
	"fmt"
	"log"
)

func main() {
	conf := config.Config{"./config/config.json"}
	conMap, err := conf.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("port: ", conMap["port"])

	//здесь запускаем сервер
}
