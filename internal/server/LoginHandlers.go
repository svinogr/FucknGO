package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"time"
)

// auth user and send jwt token
// test handler for aut from html form
func logPage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("log.html")
	tmp.Execute(w, "done")
}

// auth responses with token if log is success

func refreshTokekn() {

}

func auth(w http.ResponseWriter, r *http.Request) {
	var uM model.UserModel

	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// юзер есть с таким паролем
	validUser, err := getValidUser(uM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, _ := jwt.CreateJWTToken(validUser.Id)
	refreshToken, _ := jwt.CreateJWTRefreshToken(validUser.Id)

	// есть ли уже сесиия для данного юзера
	session, err := getSessionForUserIdIfIs(validUser.Id)
	// TODOD можно добавить проверку на срок сессии
	if err == nil {
		session.RefreshToken = refreshToken
		updateSession(session)
	} else {
		session := repo.SessionModelRepo{
			UserId:       validUser.Id,
			RefreshToken: refreshToken,
			UserAgent:    "",
			Fingerprint:  "",
			Ip:           "",
			ExpireIn:     time.Now().Add(repo.Exp_session),
			CreatedAt:    time.Now(),
		}

		createSession(session)
	}

	token := model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(token)
	if err != nil {
		log.NewLog().Fatal(err)
	}

	fmt.Fprint(w, string(marshal))
}

func updateSession(session *repo.SessionModelRepo) {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()

	_, err = sessionRepo.UpdateSession(session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func getSessionForUserIdIfIs(id uint64) (*repo.SessionModelRepo, error) {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()
	session, err := sessionRepo.FindSessionByUserId(id)

	if err != nil {

		return nil, err
	}

	return session, nil
}

func createSession(session repo.SessionModelRepo) {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()

	_, err = sessionRepo.CreateSession(&session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

// validUser gets valid user by email and password
func getValidUser(user model.UserModel) (*repo.UserModelRepo, error) {
	conf, err := config.GetConfig()

	if err != nil {
		return nil, err
	}

	userRepo := repo.NewDataBase(conf).User()

	uBemail, err := userRepo.FindUserByEmail(user.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uBemail.Password), []byte(user.Password))

	if err != nil {
		return nil, err
	}

	return uBemail, nil
}

func GetUserIdFromContext(r *http.Request) (interface{}, error) {
	value := r.Context().Value(jwt.UserId)

	if value == nil {
		return nil, errors.New("Not found id")
	}

	return value, nil
}

func logOut(w http.ResponseWriter, r *http.Request) {
	var uM model.UserModel

	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := jwt.GetClaims(uM.Token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = claims.Valid()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userId := claims[jwt.UserId]

	conf, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	repo := base.Sessions()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	_, err = repo.DeleteSessionByUserId(userId.(uint64))

	if err != nil {
		log.NewLog().Fatal(err)
	}

	http.Redirect(w, r, "/api/login", http.StatusMovedPermanently)
}
