package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	. "FucknGO/internal/jwt"
	"fmt"
	"log"
	"testing"
	"time"
)

func GetSessionRepo() (*repo.SessionRepo, error) {
	conf, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	sessionRepo := repo.NewDataBase(conf).Sessions()

	return sessionRepo, nil
}

var testUserWithToken *repo.UserModelRepo
var session *repo.SessionModelRepo
var sessionRepo *repo.SessionRepo

func CreatTestUser() {
	userRepo, err := GetUserRepo()

	if err != nil {
		log.Fatal(err)
	}

	testUser = repo.UserModelRepo{Name: "foo", Email: "emeil", Password: "pass"}

	testUserWithToken, err = userRepo.CreateUser(&testUser)

	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateSession(t *testing.T) {
	sessionRepo, err := GetSessionRepo()

	if err != nil {
		t.Error(err)
	}

	CreatTestUser()

	if err != nil {
		t.Errorf("%v", err)
	}

	refreshToken, err := CreateJWTRefreshToken(testUserWithToken.Id)

	if err != nil {
		t.Error(err)
	}

	session = &repo.SessionModelRepo{
		UserId:       0,
		RefreshToken: refreshToken,
		UserAgent:    "",
		Fingerprint:  "",
		Ip:           "",
		ExpireIn:     0,
		CreatedAt:    time.Now(),
	}

	createSession, err := sessionRepo.CreateSession(session)

	if err != nil {
		t.Error(err)
	}

	if createSession.Id < 0 {
		t.Error()
	}

	fmt.Print(session)
}

/*
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
}*/
