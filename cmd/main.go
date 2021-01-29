package main

import (
	"FucknGO/config"
	"FucknGO/log"
	"FucknGO/server"
	"flag"
	"strconv"
)

func init() {
	getConfig()
}

func main() {
	config := getConfig()

	startServer(*config)
}

func getConfig() *config.Config {
	conf, err := config.GetConfig(config.Path)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	portArg := flag.String("port", "8080", "used port")
	flag.Parse()

	if *portArg != string(conf.JsonStr.ServerConfig.Port) {

		value, err := strconv.Atoi(*portArg)

		if err == nil {
			conf.JsonStr.ServerConfig.Port = uint16(value)
		}
	}

	return conf
}

func startServer(config config.Config) {
	ser := server.Server{Config: config}
	ser.Start()
}
