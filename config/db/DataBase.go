package db

type DataBaseStr struct {
	Postgres PostgresStr `json:"postgresql"`
	Geiop    GeiopStr    `json:"geoip"`
	Redis    RedisStr    `json:"redis"`
}
