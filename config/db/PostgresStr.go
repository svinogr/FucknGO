package db

type PostgresStr struct {
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	Password string `json:"password"`
	BaseName string `json:"basename"`
	User     string `json:"user"`
}
