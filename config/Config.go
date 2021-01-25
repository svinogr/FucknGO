package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Path string
}

func (p *Config) SetPath(path string) {
	p.Path = path
}

// GetConfig return map with config from json
func (p *Config) GetConfig() (map[string]string, error) {
	var jsonConfigMap map[string]string

	fileJson, err := os.Open(p.Path)

	defer fileJson.Close()

	if err != nil {
		return jsonConfigMap, err
	}

	if err != nil {
		return jsonConfigMap, err
	}

	fileInfo, _ := os.Stat(p.Path)
	fileSize := fileInfo.Size()

	var dataJson = make([]byte, fileSize)

	_, err = fileJson.Read(dataJson)

	if err != nil {
		return jsonConfigMap, err
	}

	err = json.Unmarshal(dataJson, &jsonConfigMap)

	if err != nil {
		return jsonConfigMap, err
	}

	return jsonConfigMap, nil
}
