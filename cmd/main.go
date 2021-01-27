package main

import (
	"FucknGO/config"
	"FucknGO/server"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	setLogging()

	config := getConfig()

	startServer(*config)
}

func getConfig() *config.Config {
	conf := config.Config{}
	config, err := conf.ReadConfig(config.Path)

	if err != nil {
		log.Fatal(err)
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

func setLogging() {
	_, err := os.Stat("log.txt")

	if err != nil {
		file, err := os.Create("log.txt")

		defer file.Close()

		if err != nil {
			fmt.Print(err)
		}

		log.SetOutput(file)
	}
}

func startServer(config config.Config) {
	ser := server.Server{Config: config}
	ser.Start()
}
