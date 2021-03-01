package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/db/user"
	. "FucknGO/internal/jwt"
	"database/sql"
	"log"
	"testing"
)

func GetTokenRepo() (*repo.TokenRepo, error) {
	conf, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	tokenRepo := repo.NewDataBase(conf).Token()

	return tokenRepo, nil
}

var testUserWithToken *user.UserModelRepo
var token *user.TokenModelRepo
var tokenRepo *repo.TokenRepo

func CreatTestUser() {
	userRepo, err := GetUserRepo()

	if err != nil {
		log.Fatal(err)
	}

	testUser = user.UserModelRepo{Name: "foo", Email: "emeil", Password: "pass"}

	testUserWithToken, err = userRepo.CreateUser(&testUser)

	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateToken(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Error(err)
	}

	CreatTestUser()

	if err != nil {
		t.Errorf("%v", err)
	}

	createJWT, err := CreateJWT(1)

	if err != nil {
		t.Error(err)
	}

	token = &user.TokenModelRepo{
		Token:  createJWT,
		UserId: testUserWithToken.Id,
	}

	createToken, err := tokenRepo.CreateToken(token)

	if err != nil {
		t.Error(err)
	}

	if createToken.Id < 0 {
		t.Error()
	}
}

func TestFindTokenByUserId(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Error(err)
	}

	findToken, err := tokenRepo.FindTokenByUserId(testUserWithToken.Id)

	if err != nil {
		t.Error(err)
	}

	if findToken.Id < 0 {
		t.Error()
	}
}

func TestFindTokenByIdNotAddedInDB(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Error(err)
	}

	_, err = tokenRepo.FindTokenByUserId(testUserWithToken.Id + 1)

	if err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestUpdateToken(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Error(err)
	}

	newToken := "newToken"
	token.Token = newToken

	updateToken, err := tokenRepo.UpdateToken(token)

	if err != nil {
		t.Error(err)
	}

	if updateToken.Id != token.Id && updateToken.Token != newToken && updateToken.UserId != token.UserId {
		t.Error()
	}
}

func TestUpdateTokenNotAddedInDB(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Error(err)
	}

	newToken := "newToken"
	token.Token = newToken
	token.UserId = testUserWithToken.Id + 1

	_, err = tokenRepo.UpdateToken(token)

	if err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestDeleteToken(t *testing.T) {
	tokenRepo, err := GetTokenRepo()

	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = tokenRepo.DeleteToken(token)

	if err != nil {
		t.Error(err)
	}
}
