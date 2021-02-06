package db

import (
	"time"
)

type Slave struct {
	ID        uint
	CreatedAt time.Time
	Port      int
	StaticDir string
}

func Create(slave *Slave) error {
	return db.Create(slave).Error
}

func GetList() (slaves []Slave, err error) {
	if err = db.Find(&slaves).Error; err != nil {
		return nil, err
	}
	return slaves, err
}
func GetByID(id uint) (slave Slave, err error) {
	slave.ID = id
	err = db.First(&slave).Error
	return slave, err
}

func Delete(id uint) error {
	return db.Delete(&Slave{ID: id}).Error
}
