package repo

import "errors"

const (
	TABLE_NAME_SHOP   = "shops"
	COL_SHOP_ID       = "id"
	COL_SHOP_COORD_ID = "coord_id"
	COL_SHOP_NAME     = "name"
	COL_SHOP_ADDRESS  = "address"
	COL_SHOP_USER_ID  = "user_id"
)

type ShopModelRepo struct {
	Id      uint64
	UserId  uint64
	CoordId uint64
	Name    string
	Address string
}

type ShopRepo struct {
	db *DataBase
}

func (c *ShopRepo) Create(shop *ShopModelRepo) (*ShopModelRepo, error) {
	if err := c.db.Db.QueryRow("INSERT into "+TABLE_NAME_SHOP+
		" ("+COL_SHOP_COORD_ID+", "+COL_SHOP_USER_ID+", "+COL_SHOP_NAME+", "+
		COL_SHOP_ADDRESS+") VALUES ($1, $2, $3, $4) RETURNING "+COL_SHOP_ID,
		shop.CoordId,
		shop.UserId,
		shop.Name,
		shop.Address).
		Scan(&shop.Id); err != nil {
		return nil, err
	}

	return shop, nil
}

func (c *ShopRepo) Update(shop *ShopModelRepo) (*ShopModelRepo, error) {
	err := c.db.Db.QueryRow("UPDATE "+TABLE_NAME_SHOP+" set "+
		COL_SHOP_COORD_ID+"= $1, "+
		COL_SHOP_USER_ID+"=$2,"+
		COL_SHOP_NAME+"= $3, "+
		COL_SHOP_ADDRESS+"= $4"+
		" WHERE "+COL_SHOP_ID+"=$5 returning "+
		COL_SHOP_ID+", "+
		COL_SHOP_USER_ID+", "+
		COL_SHOP_COORD_ID+", "+
		COL_SHOP_NAME+", "+
		COL_SHOP_ADDRESS,
		shop.CoordId,
		shop.UserId,
		shop.Name,
		shop.Address,
		shop.Id).
		Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address)

	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (c *ShopRepo) Delete(shop *ShopModelRepo) (*ShopModelRepo, error) {
	_, err := c.db.Db.Exec("DELETE from "+TABLE_NAME_SHOP+" where "+COL_SHOP_ID+" = $1", shop.Id)

	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (c *ShopRepo) FindByTitle(shop *ShopModelRepo) (*ShopModelRepo, error) {
	if err := c.db.Db.QueryRow("SELECT * from "+TABLE_NAME_SHOP+" where "+COL_SHOP_ID+"=$1",
		shop.Id).Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address); err != nil {

		return nil, err
	}

	return shop, nil
}

func (c *ShopRepo) FindAll() (*[]ShopModelRepo, error) {
	shopList := []ShopModelRepo{}

	row, err := c.db.Db.Query("SELECT * from " + TABLE_NAME_SHOP)

	if err != nil {
		return &shopList, err
	}

	for row.Next() {
		shop := ShopModelRepo{}
		err := row.Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address)

		if err != nil {
			return nil, err
		}
		shopList = append(shopList, shop)
	}

	return &shopList, nil
}

// FindByListCoords finds all shoops by id coord
func (c *ShopRepo) FindByListCoords(coords *[]CoordModelRepo) (*[]ShopModelRepo, error) {
	shopsList := []ShopModelRepo{}

	for _, ell := range *coords {
		shop := ShopModelRepo{}

		row := c.db.Db.QueryRow("SELECT * from "+TABLE_NAME_SHOP+" where "+COL_SHOP_ID+"=$1", ell.Id)

		err := row.Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address)

		if err == errors.New("ErrNoRows") {
			continue
		}

		if err != nil {
			return nil, err
		}

		shopsList = append(shopsList, shop)
	}

	return &shopsList, nil
}

func (c *ShopRepo) FindByUserId(user UserModelRepo) (*[]ShopModelRepo, error) {
	shopList := []ShopModelRepo{}

	row, err := c.db.Db.Query("SELECT * from "+TABLE_NAME_SHOP+" where "+COL_SHOP_USER_ID+"=$1", user.Id)

	if err != nil {
		return &shopList, err
	}

	for row.Next() {
		shop := ShopModelRepo{}
		err := row.Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address)

		if err != nil {
			return nil, err
		}
		shopList = append(shopList, shop)
	}

	return &shopList, nil
}

func (c *ShopRepo) FindById(id uint64) (*ShopModelRepo, error) {
	shop := ShopModelRepo{}

	if err := c.db.Db.QueryRow("SELECT * from "+TABLE_NAME_SHOP+" where "+COL_SHOP_ID+"=$1",
		id).Scan(&shop.Id, &shop.UserId, &shop.CoordId, &shop.Name, &shop.Address); err != nil {

		return nil, err
	}

	return &shop, nil

}
