package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port int8
}

func main() {
	fileJson, err := os.Open("./config/config.json")

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	defer fileJson.Close()

	dataJson := make([]byte, 64)

	_, err = fileJson.Read(dataJson)

	fmt.Println(dataJson)

	var jsonConfig map[string] string

	json.Unmarshal(dataJson, jsonConfig)

	fmt.Print(jsonConfig["port"]) // вот тут нил


}
