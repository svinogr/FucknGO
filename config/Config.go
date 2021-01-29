package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Config struct {
	JsonStr jsonStr
}

var instance *Config
var once sync.Once

func GetConfig(path string) (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}
		instance, err = instance.readConfig(path)
	})

	return instance, err
}

func (p *Config) readConfig(path string) (*Config, error) {
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
