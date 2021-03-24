package json

import (
	"FucknGO/config/db"
	"FucknGO/config/log"
	"FucknGO/config/rabbitmq"
	"FucknGO/config/server"
	"FucknGO/config/ui"
)

const Path = "./config/json/config.json"

type JsonStr struct {
	ServerConfig server.ServerConfigStr `json:"server"`
	UiConfig     ui.UiConfigStr         `json:"ui"`
	Log          log.LogStr             `json:"log"`
	DataBase     db.DataBaseStr         `json:"databases"`
	RabbitMQ     rabbitmq.RabbitMQStr   `json:"rabbit_mq"`
}
