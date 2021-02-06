package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var conf *Config

type Config struct {
	Basic        *Basic
	MasterServer *MasterServer
}

type Basic struct {
	Debug bool `json:"debug"`
}

type MasterServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func read(name string, template interface{}) (interface{}, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(fmt.Sprintf("./configs/%s.json", name))
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("Successfully Opened %s.json\n", name)
	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// read our opened jsonFile as a byte array.
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.Unmarshal(bytes, &template)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return template, nil
}

func checkError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func make() *Config {
	basic, err := read("basic", &Basic{})
	master, err := read("master", &MasterServer{})

	conf = &Config{
		Basic:        basic.(*Basic),
		MasterServer: master.(*MasterServer),
	}

	if mPort := os.Getenv("MASTER_PORT"); mPort != "" {
		conf.MasterServer.Port, err = strconv.Atoi(mPort)
	}

	flag.IntVar(&conf.MasterServer.Port, "master_port", conf.MasterServer.Port,
		"bind `port` for main HTTP server")
	flag.Parse()

	checkError(err)

	return conf
}

func Get() *Config {
	if conf != nil {
		return conf
	} else {
		return make()
	}
}
