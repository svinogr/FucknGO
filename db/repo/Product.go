package repo

const (
	TABLE_NAME_PRODUCT        = "product"
	COL_PRODUCT_ID            = "id"
	COL_PRODUCT_SHOP_STOCK_ID = "shop_stock_id"
	COL_PRODUCT_TITLE         = "title"
	COL_PRODUCT_DESCRIPTION   = "description"
	COL_PRODUCT_PRICE         = "price"
	COL_PRODUCT_PRICE_OLD     = "price_old"
)

type ProductModelRepo struct {
	Id          uint64
	ShopStockId uint64
	Title       string
	Description string
	Price       float64
	PriceOld    float64
}

type ProductRepo struct {
	db *DataBase
}

func (p *ProductRepo) Create(product *ProductModelRepo) (*ProductModelRepo, error) {
	if err := p.db.Db.QueryRow("INSERT into "+TABLE_NAME_PRODUCT+
		" ("+COL_PRODUCT_SHOP_STOCK_ID+", "+COL_PRODUCT_TITLE+", "+
		""+COL_PRODUCT_DESCRIPTION+", "+COL_PRODUCT_PRICE+", "+
		COL_PRODUCT_PRICE_OLD+") VALUES ($1, $2, $3, $4, $5, $6) RETURNING "+COL_PRODUCT_ID,
		product.ShopStockId,
		product.Title,
		product.Description,
		product.Price,
		product.PriceOld).
		Scan(&product.Id); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepo) Update(product *ProductModelRepo) (*ProductModelRepo, error) {
	err := p.db.Db.QueryRow("UPDATE "+TABLE_NAME_PRODUCT+" set "+
		COL_PRODUCT_SHOP_STOCK_ID+"= $1, "+
		COL_PRODUCT_TITLE+"= $2, "+
		COL_PRODUCT_DESCRIPTION+"= $3"+
		COL_PRODUCT_PRICE+"= $4"+
		COL_PRODUCT_PRICE_OLD+"= $5"+
		" WHERE "+COL_PRODUCT_ID+"=$6 returning "+COL_PRODUCT_ID+", "+
		COL_PRODUCT_SHOP_STOCK_ID+", "+
		COL_PRODUCT_TITLE+", "+
		COL_PRODUCT_DESCRIPTION+", "+
		COL_PRODUCT_PRICE+", "+
		COL_PRODUCT_PRICE_OLD,
		product.ShopStockId,
		product.Title,
		product.Description,
		product.Price,
		product.PriceOld,
		product.Id).
		Scan(&product.Id, &product.ShopStockId, &product.Title, &product.Description, &product.Price, &product.PriceOld)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepo) Delete(product *ProductModelRepo) (*ProductModelRepo, error) {
	_, err := p.db.Db.Exec("DELETE from "+TABLE_NAME_PRODUCT+" where "+COL_PRODUCT_ID+" = $1", product.Id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepo) FindByTitle(product *ProductModelRepo) (*ProductModelRepo, error) {
	if err := p.db.Db.QueryRow("SELECT * from "+TABLE_NAME_SHOP+" where "+COL_PRODUCT_TITLE+"=$1",
		product.Title).Scan(&product.Id, &product.Title, &product.Description, &product.Price, &product.PriceOld); err != nil {

		return nil, err
	}

	return product, nil
}

func (p *ProductRepo) FindAll() (*[]ProductModelRepo, error) {
	productList := []ProductModelRepo{}

	row, err := p.db.Db.Query("SELECT * from " + TABLE_NAME_PRODUCT)

	if err != nil {
		return &productList, err
	}

	for row.Next() {
		product := ProductModelRepo{}
		err := row.Scan(&product.Id, &product.ShopStockId, &product.Title, &product.Description, &product.Price, &product.PriceOld)

		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	return &productList, nil
}
