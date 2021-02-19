package db

import (
	"FucknGO/config"
	"FucknGO/db/repo"
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
	Db          sql.DB
	userRepo    *repo.UserRepo
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

	db, err := sql.Open("postgres", urlDB)

	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *DataBase) CloseDataBase() error {
	return d.db.Close()
}

func (d *DataBase) User() *repo.UserRepo {
	if d.userRepo != nil {
		return d.userRepo
	}

	d.userRepo = &repo.UserRepo{
		d,
	}

	return d.userRepo
}
