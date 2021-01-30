package config

const Path = "./config/jsondata/config.json"

type jsonStr struct {
	ServerConfig ServerConfigStr `json:"server"`
	UiConfig     UiConfigStr     `json:"ui"`
	Log          LogStr          `json:"log"`
}
