package config

import (
	js "FucknGO/config/json"
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	JsonStr js.JsonStr
}

var instance *Config
var once sync.Once

// GetConfig return reading config from json
func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}
		instance, err = instance.readConfig()
	})

	return instance, err
}

// readConfig read config from json file defined in value js.Path
func (p *Config) readConfig() (*Config, error) {
	path := js.Path
	fileJson, err := os.Open(path)

	defer fileJson.Close()

	if err != nil {
		return nil, err
	}

	fileInfo, _ := os.Stat(path)
	fileSize := fileInfo.Size()

	var dataJson = make([]byte, fileSize)

	_, err = fileJson.Read(dataJson)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(dataJson, &p.JsonStr)

	if err != nil {
		return nil, err
	}

	return p, nil
}
