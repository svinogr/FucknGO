package server

type ServerConfigStr struct {
	Port              uint16 `json:"port"`
	Ws                string `json:"ws"`
	MiddlewareTimeout string `json:"middlewareTimeout"`
	JwtSecret         string `json:"jwtSecret"`
	JwtLifeTimeDays   uint16 `json:"jwtLifeTimeDays"`
}

//TODO не забыть использовтаь ключ для jwt отсюда
