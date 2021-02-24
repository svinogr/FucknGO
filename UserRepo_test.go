package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/db/user"
	"testing"
)

func GetDB() (*repo.DataBase, error) {
	conf, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	return repo.NewDataBase(conf), err
}

var u = user.UserModelRepo{Name: "foo", Email: "1@mail.test", Password: "123456"}

func TestCreateUser(t *testing.T) {
	db, err := GetDB()

	if err != nil {
		t.Errorf("%v", err)
	}

	userRepo := db.User()
	userRepo.CreateUser(&u)

}
