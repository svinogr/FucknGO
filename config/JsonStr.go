package config

const Path = "./config/jsondata/config.json"

type JsonStr struct {
	ServerConfig ServerConfigStr `json:"server"`
	UiConfig     UiConfigStr     `json:"ui"`
}
