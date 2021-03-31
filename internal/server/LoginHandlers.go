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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// юзер есть с таким паролем
	validUser, err := getValidUser(uM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenAccess, _ := jwt.CreateJWTToken(validUser.Id)
	tokenRefresh, _ := jwt.CreateJWTRefreshToken(validUser.Id)

	// есть ли уже сесиия для данного юзера
	ok := hasSessionForUserId(validUser.Id)

	if ok {
		session := repo.SessionModelRepo{
			UserId:       validUser.Id,
			RefreshToken: tokenRefresh,
			UserAgent:    "",
			Fingerprint:  "",
			Ip:           "",
			ExpireIn:     0,
		}

		updateSession(session)
	} else {
		createSession(tokenRefresh)
	}

	token := model.TokenModel{
		AccessToken:  tokenAccess,
		RefreshToken: tokenRefresh,
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, token)
}

func updateSession(session repo.SessionModelRepo) {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()

	_, err = sessionRepo.UpdateSession(&session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func hasSessionForUserId(id uint64) bool {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()
	session, err := sessionRepo.FindSessionByUserId(id)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	if session.Id > 0 {
		return true
	}

	return false
}

func createSession(refreshToken string) {
	conf, err := config.GetConfig()
	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	sessionRepo := base.Sessions()

	session := repo.SessionModelRepo{
		UserId:       0,
		RefreshToken: refreshToken,
		UserAgent:    "",
		Fingerprint:  "",
		Ip:           "",
		ExpireIn:     0,
		CreatedAt:    time.Now(),
	}

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

	//userId := claims[jwt.UserId]
	var userId uint64 = 24

	conf, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	repo := base.Token()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	// err = repo.DeleteTokenByUserId(userId.(uint64))
	_, err = repo.DeleteTokenByUserId(userId)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	fmt.Fprint(w, nil)
}
