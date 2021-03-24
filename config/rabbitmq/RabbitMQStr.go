package rabbitmq

type RabbitMQStr struct {
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	Password string `json:"password"`
	User     string `json:"user"`
}
