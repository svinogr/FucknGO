package main

import (
	"FucknGO/config"
	"FucknGO/config/json"
	"FucknGO/db"
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
	//postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]

	//fmt.Printf("postgresql://[%s[:%s]@][%s][:%d][/postgres]","postgres", "postgres","localhost", 5432)

	db := db.NewDataBase(conf)
	err := db.OpenDataBase()
	fmt.Print(err)
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
