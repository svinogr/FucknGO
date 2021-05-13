package model

type ShopModel struct {
	Id       uint64       `json:"id"`
	UserId   uint64       `json:"user_id"`
	CoordLat string       `json:"coord_lat"`
	CoordLng string       `json:"coord_lng"`
	Name     string       `json:"name"`
	Address  string       `json:"address"`
	Stocks   []StockModel `json:"stoks"`
}
