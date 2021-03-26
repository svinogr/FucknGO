package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
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

var testUser = repo.UserModelRepo{Name: "foo", Email: "email", Password: "pass"}

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

	newUser := repo.UserModelRepo{
		Id:       testUser.Id,
		Name:     "newFoo",
		Password: "newPass",
		Email:    "newEmail",
	}

	testUser, err := userRepo.UpdateUser(&newUser)

	if err != nil {
		t.Error(err)
	}

	if newUser.Id != testUser.Id && newUser.Name != testUser.Name && newUser.Password != testUser.Password &&
		newUser.Email != testUser.Email {
		t.Error()
	}
}

func TestUpdateUserNotAddedInDB(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	newUser := repo.UserModelRepo{
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

func TestFindUserByEmail(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	user, err := userRepo.FindUserByEmail(testUser.Email)

	if err != nil {
		t.Error(err)
	}

	if user.Id != testUser.Id && user.Name != testUser.Name && user.Password != testUser.Password &&
		user.Email != testUser.Email {
		t.Error()
	}

}

func TestFindUserByName(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	user, err := userRepo.FindUserByName(testUser.Name)

	if err != nil {
		t.Error(err)
	}

	if user.Id != testUser.Id && user.Name != testUser.Name && user.Password != testUser.Password &&
		user.Email != testUser.Email {
		t.Error()
	}

}

func TestFindAllUser(t *testing.T) {
	userRepo, err := GetUserRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	userList, err := userRepo.FindAllUser()

	if err != nil {
		t.Error(err)
	}

	if userList != nil {
		if len(*userList) < 1 {
			t.Error()
		}
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
