package repo

import (
	"FucknGO/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var driverSql = "postgres"
var dbUrlStart = "postgresql"
var sslMode = "sslmode=disable"

type DataBase struct {
	config      *config.Config
	dataBaseUrl string
	Db          *sql.DB
	userRepo    *UserRepo
	tokenRepo   *TokenRepo
}

func NewDataBase(config *config.Config) *DataBase {
	return &DataBase{
		config: config,
	}
}

func (d *DataBase) OpenDataBase() error {
	urlDB := fmt.Sprintf("%s://%s:%s@%s:%d/%s?%s",
		dbUrlStart,
		d.config.JsonStr.DataBase.Postgres.User,
		d.config.JsonStr.DataBase.Postgres.Password,
		d.config.JsonStr.DataBase.Postgres.Address,
		d.config.JsonStr.DataBase.Postgres.Port,
		d.config.JsonStr.DataBase.Postgres.BaseName,
		sslMode)

	var err error
	d.Db, err = sql.Open("postgres", urlDB)

	if err != nil {
		return err
	}

	if err = d.Db.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *DataBase) CloseDataBase() error {
	return d.Db.Close()
}

func (d *DataBase) User() *UserRepo {
	if d.userRepo != nil {
		return d.userRepo
	}
	d.OpenDataBase()
	d.userRepo = &UserRepo{
		d,
	}

	return d.userRepo
}

func (d *DataBase) Token() *TokenRepo {
	if d.tokenRepo != nil {
		return d.tokenRepo
	}

	d.tokenRepo = &TokenRepo{
		d,
	}

	return d.tokenRepo
}
