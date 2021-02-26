package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/db/user"
	"database/sql"
	"testing"
)

func GetUserRepo() (*repo.UserRepo, error) {
	conf, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	userRepo := repo.NewDataBase(conf).User()

	return userRepo, nil
}

var testUser = user.UserModelRepo{Name: "foo", Email: "emeil", Password: "pass"}

func TestCreateUser(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	createUser, err := userRepo.CreateUser(&testUser)

	if err != nil {
		t.Error(err)
	}

	if createUser.Id < 0 {
		t.Error()
	}
}

func TestFindUserById(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	findUser, err := userRepo.FindUserById(testUser.Id)

	if err != nil {
		t.Error(err)
	}

	if findUser.Id < 0 {
		t.Error()
	}
}

func TestFindUserByIdNotAddedInDB(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = userRepo.FindUserById(testUser.Id + 1)

	if err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestUpdateUser(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	newUser := user.UserModelRepo{
		Id:       testUser.Id,
		Name:     "newFoo",
		Password: "newPass",
		Email:    "newEmail",
	}

	updateUser, err := userRepo.UpdateUser(&newUser)

	if err != nil {
		t.Error(err)
	}

	if updateUser.Id != testUser.Id && updateUser.Name != testUser.Name && updateUser.Password != testUser.Password &&
		updateUser.Email != testUser.Email {
		t.Error()
	}
}

func TestUpdateUserNotAddedInDB(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	newUser := user.UserModelRepo{
		Id:       testUser.Id + 1,
		Name:     "newFoo",
		Password: "newPass",
		Email:    "newEmail",
	}

	_, err = userRepo.UpdateUser(&newUser)

	if err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = userRepo.DeleteUser(&testUser)

	if err != nil {
		t.Error(err)
	}
}
