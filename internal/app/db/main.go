package db

import (
	"FucknGO/internal/app/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB //база данных
var err error

func connect() {
	conf := config.Get()

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		conf.DB.Host, conf.DB.User, conf.DB.Dbname, conf.DB.Password)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	checkError(err)

	// Migrate the schema
	err = db.AutoMigrate(&Slave{})
	checkError(err)
}

// возвращает дескриптор объекта DB
func GetDB() *gorm.DB {
	if db == nil {
		connect()
	}
	return db
}

func checkError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
