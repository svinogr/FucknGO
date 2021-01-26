package config

type ServerConfigStr struct {
	Port              uint16 `json:"port"`
	Ws                string `json:"ws"`
	MiddlewareTimeout string `json:"middlewareTimeout"`
	JwtSecret         string `json:"jwtSecret"`
	JwtLifeTimeDays   int8   `json:"jwtLifeTimeDays"`
}
