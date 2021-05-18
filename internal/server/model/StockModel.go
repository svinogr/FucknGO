package model

type StockModel struct {
	Id          uint64 `json:"id"`
	ShopId      uint64 `json:"shop_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateStart   string `json:"date_start"`
	DateFinish  string `json:"date_finish"`
}
