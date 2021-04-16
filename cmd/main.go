package main

import (
	"FucknGO/config"
	"FucknGO/internal/server"
	"FucknGO/log"
	"flag"
	"fmt"
	"strconv"
)

func init() {
	conf = *getConfig()
}

var conf config.Config

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}()

	startServer(conf)
}

func getConfig() *config.Config {
	conf, err := config.GetConfig()

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
	port := fmt.Sprint(config.JsonStr.ServerConfig.Port)

	fb, err := server.FabricServer()

	if err != nil {
		//log.NewLog().Fatal(err)
	}

	ser := fb.GetNewMasterServer("0.0.0.0", port)

	err = fb.RunServer(ser)

	if err.Error() != "http: serverApi closed" {
		panic(err)
		log.NewLog().Fatal(err)
	}

	log.NewLog().PrintError(err)
}
