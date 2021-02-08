package db

import (
	"FucknGO/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var driverSql = "postgres"
var dbUrlStart = "postgresql"
var sslMode = "sslmode=disable"

type Database struct {
	config      *config.Config
	dataBaseUrl string
	db          sql.DB
}

func NewDataBase(config *config.Config) *Database {
	return &Database{
		config: config,
	}
}

func (d *Database) OpenDataBase() error {
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

func (d *Database) CloseDataBase() error {
	return d.db.Close()
}
