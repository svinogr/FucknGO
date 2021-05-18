package repo

import "time"

const (
	TABLE_NAME_STOCK           = "shops_stock"
	COL_STOCK_ID               = "id"
	COL_STOCK_SHOP_ID          = "shop_id"
	COL_STOCK_SHOP_TITLE       = "title"
	COL_STOCK_SHOP_DESCRIPTON  = "description"
	COL_STOCK_SHOP_DATE_START  = "date_start"
	COL_STOCK_SHOP_DATE_FINISH = "date_finish"
)

type ShopStockModelRepo struct {
	Id          uint64
	ShopId      uint64
	Title       string
	Description string
	DateStart   time.Time
	DateFinish  time.Time
}

type ShopStockRepo struct {
	db *DataBase
}

func (s *ShopStockRepo) Create(shopStock *ShopStockModelRepo) (*ShopStockModelRepo, error) {
	if err := s.db.Db.QueryRow("INSERT into "+TABLE_NAME_STOCK+
		" ("+COL_STOCK_SHOP_ID+", "+
		COL_STOCK_SHOP_TITLE+", "+
		COL_STOCK_SHOP_DESCRIPTON+", "+
		COL_STOCK_SHOP_DATE_START+", "+
		COL_STOCK_SHOP_DATE_FINISH+") VALUES ($1, $2, $3, $4, $5) RETURNING "+COL_STOCK_ID,
		shopStock.ShopId,
		shopStock.Title,
		shopStock.Description,
		shopStock.DateStart,
		shopStock.DateFinish).
		Scan(&shopStock.Id); err != nil {
		return nil, err
	}

	return shopStock, nil
}

func (s *ShopStockRepo) Update(stock *ShopStockModelRepo) (*ShopStockModelRepo, error) {
	err := s.db.Db.QueryRow("UPDATE "+TABLE_NAME_STOCK+" set "+
		COL_STOCK_SHOP_ID+"= $1, "+
		COL_STOCK_SHOP_TITLE+"= $2, "+
		COL_STOCK_SHOP_DESCRIPTON+"= $3, "+
		COL_STOCK_SHOP_DATE_START+"= $4, "+
		COL_STOCK_SHOP_DATE_FINISH+"= $5 "+
		"where "+COL_STOCK_ID+"=$6 returning "+
		COL_STOCK_ID+", "+
		COL_STOCK_SHOP_ID+","+
		COL_STOCK_SHOP_TITLE+", "+
		COL_STOCK_SHOP_DESCRIPTON+", "+
		COL_STOCK_SHOP_DATE_START+", "+
		COL_STOCK_SHOP_DATE_FINISH,
		stock.ShopId,
		stock.Title,
		stock.Description,
		stock.DateStart,
		stock.DateFinish,
		stock.Id).
		Scan(&stock.Id, &stock.ShopId, &stock.Title, &stock.Description, &stock.DateStart, &stock.DateFinish)

	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *ShopStockRepo) Delete(shopStock *ShopStockModelRepo) (*ShopStockModelRepo, error) {
	_, err := s.db.Db.Exec("DELETE from "+TABLE_NAME_STOCK+" where "+COL_STOCK_ID+" = $1", shopStock.Id)

	if err != nil {
		return nil, err
	}

	return shopStock, nil
}

func (s *ShopStockRepo) FindById(shopStock *ShopStockModelRepo) (*ShopStockModelRepo, error) {
	if err := s.db.Db.QueryRow("SELECT * from "+TABLE_NAME_STOCK+" where "+COL_STOCK_ID+"=$1",
		shopStock.Id).
		Scan(&shopStock.Id, &shopStock.ShopId, &shopStock.Title, &shopStock.Description, &shopStock.DateStart, &shopStock.DateFinish); err != nil {

		return nil, err
	}

	return shopStock, nil
}

// FindByShop find all stocks by shop
func (s *ShopStockRepo) FindByShop(shop *ShopModelRepo) (*[]ShopStockModelRepo, error) {
	shopStockList := []ShopStockModelRepo{}

	row, err := s.db.Db.Query("SELECT * from "+TABLE_NAME_STOCK+" where "+COL_STOCK_SHOP_ID+"=$1",
		shop.Id)

	if err != nil {
		return &shopStockList, err
	}

	for row.Next() {
		shopStock := ShopStockModelRepo{}
		err := row.Scan(&shopStock.Id, &shopStock.ShopId, &shopStock.Title, &shopStock.Description, &shopStock.DateStart, &shopStock.DateFinish)

		if err != nil {
			return nil, err
		}

		shopStockList = append(shopStockList, shopStock)
	}

	return &shopStockList, nil
}

func (s *ShopStockRepo) FindAll() (*[]ShopStockModelRepo, error) {
	shopStockList := []ShopStockModelRepo{}

	row, err := s.db.Db.Query("SELECT * from " + TABLE_NAME_STOCK)

	if err != nil {
		return &shopStockList, err
	}

	for row.Next() {
		shopStock := ShopStockModelRepo{}
		err := row.Scan(&shopStock.Id, &shopStock.ShopId, &shopStock.Title, &shopStock.Description, &shopStock.DateStart, &shopStock.DateFinish)

		if err != nil {
			return nil, err
		}
		shopStockList = append(shopStockList, shopStock)
	}

	return &shopStockList, nil
}
