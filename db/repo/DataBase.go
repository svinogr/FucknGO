package repo

import (
	"FucknGO/config"
	"FucknGO/log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	DRIVER_SQL   = "postgres"
	DB_URL_START = "postgresql"
	SSL_MODE     = "sslmode=disable"
)

type DataBase struct {
	config      *config.Config
	dataBaseUrl string
	Db          *sql.DB
	userRepo    *UserRepo
	sessionRepo *SessionRepo
}

// get new db or return created db
/*func NewDataBase(config *config.Config) *DataBase {
	return &DataBase{
		config: config,
	}
}*/

func NewDataBaseWithConfig() *DataBase {
	config, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	return &DataBase{
		config: config,
	}
}

func (d *DataBase) OpenDataBase() (err error) {
	urlDB := fmt.Sprintf("%s://%s:%s@%s:%d/%s?%s",
		DB_URL_START,
		d.config.JsonStr.DataBase.Postgres.User,
		d.config.JsonStr.DataBase.Postgres.Password,
		d.config.JsonStr.DataBase.Postgres.Address,
		d.config.JsonStr.DataBase.Postgres.Port,
		d.config.JsonStr.DataBase.Postgres.BaseName,
		SSL_MODE)

	if d.Db, err = sql.Open(DRIVER_SQL, urlDB); err != nil {
		return err
	}

	if err = d.Db.Ping(); err != nil {
		return err
	}

	return err
}

func (d *DataBase) CloseDataBase() error {
	return d.Db.Close()
}

// get UserRepository
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

//Get SessionRepository
func (d *DataBase) Sessions() *SessionRepo {
	if d.sessionRepo != nil {
		return d.sessionRepo
	}
	d.OpenDataBase()
	d.sessionRepo = &SessionRepo{
		d,
	}

	return d.sessionRepo
}
