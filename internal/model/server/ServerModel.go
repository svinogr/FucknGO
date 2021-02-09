package server

type ServerModel struct {
	Id             uint64 `json:"id"`
	StaticResource string `json:"static_resource"`
	Port           string `json:"port"`
	Address        string `json:"address"`
}
