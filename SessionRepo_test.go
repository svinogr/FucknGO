package FucknGO

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	. "FucknGO/internal/jwt"
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
		UserId:       testUserWithToken.Id,
		RefreshToken: refreshToken,
		UserAgent:    "pc",
		Fingerprint:  "hrome",
		Ip:           "2.2.2.2",
		ExpireIn:     time.Now().Add(time.Minute),
		CreatedAt:    time.Now(),
	}

	createSession, err := sessionRepo.CreateSession(session)

	if err != nil {
		t.Error(err)
	}

	if createSession.Id < 0 {
		t.Error()
	}
}

func TestFindSessionByUserId(t *testing.T) {
	sessionRepo, err := GetSessionRepo()

	if err != nil {
		t.Error(err)
	}

	findSession, err := sessionRepo.FindSessionByUserId(testUserWithToken.Id)

	if err != nil {
		t.Error(err)
	}

	if findSession.Id < 0 {
		t.Error()
	}
}

func TestUpdateSession(t *testing.T) {
	sessionRepo, err := GetSessionRepo()

	if err != nil {
		t.Error(err)
	}

	session.RefreshToken = "newToken"
	session.UserAgent = "device"
	session.Fingerprint = "browser"
	session.Ip = "1.1.1.1"

	updateSession, err := sessionRepo.UpdateSession(session)

	if err != nil {
		t.Error(err)
	}

	if updateSession.Id != session.Id && updateSession.RefreshToken != session.RefreshToken &&
		updateSession.UserAgent != session.UserAgent && updateSession.Fingerprint != session.Fingerprint &&
		updateSession.Ip != session.Ip {
		t.Error()
	}
}

func TestDeleteToken(t *testing.T) {
	sessionRepo, err := GetSessionRepo()

	if err != nil {
		t.Error(err)
	}

	rows, err := sessionRepo.DeleteSessionByUserId(session.UserId)

	if err != nil {
		t.Error(err)
	}

	if rows < 1 {
		t.Error()
	}
}
