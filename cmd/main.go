package main

import (
	"FucknGO/config"
	"FucknGO/log"
	"FucknGO/server"
	"flag"
	"strconv"
)

func main() {
	log.NewLog().PrintCommon("ddfd")
	config := getConfig()

	startServer(*config)
}

func getConfig() *config.Config {
	conf := config.Config{}
	config, err := conf.ReadConfig(config.Path)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	portArg := flag.String("port", "8080", "used port")
	flag.Parse()

	if *portArg != string(config.JsonStr.ServerConfig.Port) {

		value, err := strconv.Atoi(*portArg)

		if err == nil {
			config.JsonStr.ServerConfig.Port = uint16(value)
		}
	}

	return config
}

func startServer(config config.Config) {
	ser := server.Server{Config: config}
	ser.Start()
}
