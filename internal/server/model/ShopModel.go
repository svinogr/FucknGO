package model

type ShopModel struct {
	Id      uint64 `json:"id"`
	CoordId uint64 `json:"coordId"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
