package main

import (
	"FucknGO/config"
	"FucknGO/config/json"
	"FucknGO/log"
	"FucknGO/server"
	"flag"
	"fmt"
	"strconv"
)

func init() {
	conf = *getConfig()
}

var conf config.Config

func main() {
	fmt.Printf("new server %d", conf.JsonStr.ServerConfig.Port)
	startServer(conf)
}

func getConfig() *config.Config {
	conf, err := config.GetConfig(json.Path)

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
