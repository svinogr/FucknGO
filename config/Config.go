package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	JsonStr JsonStr
}

func (p *Config) ReadConfig(path string) (*Config, error) {
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
