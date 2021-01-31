package db

import (
	"FucknGO/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var driverSql = "postgres"

type Database struct {
	config      config.Config
	dataBaseUrl string
	db          sql.DB
}

func NewDataBase(config config.Config) Database {
	return Database{
		config: config,
	}
}

func (d *Database) OpenDataBase() error {
	db, err := sql.Open(driverSql, "postgresql://[postgres[:postgres]@][localhost][:5432][/postgres]")

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
