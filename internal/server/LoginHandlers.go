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
		session.ExpireIn = time.Now().Add(repo.Exp_session)
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

	//TODO проверить есть ли уже сессия ?? рабоатет без проверки
	err = deleteSession(uint64(userId.(float64)))

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func deleteSession(userId uint64) error {
	conf, err := config.GetConfig()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	base := repo.NewDataBase(conf)

	repo := base.Sessions()

	_, err = repo.DeleteSessionByUserId(userId)

	return err

}

//auth/refresh-tokens
func refreshToken(w http.ResponseWriter, r *http.Request) {
	tM := model.TokenModel{}
	if err := json.NewDecoder(r.Body).Decode(&tM); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshToken := tM.RefreshToken

	claims, err := jwt.GetClaims(refreshToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId := uint64(claims[jwt.UserId].(float64))
	expToken := time.Unix(int64(claims[jwt.ExpToken].(float64)), 0)

	// проверяем если срок токена меньше данного момента то
	if expToken.Before(time.Now()) {
		http.Error(w, errors.New("Token is expired").Error(), http.StatusUnauthorized)
		return
	}

	session, err := getSessionForUserIdIfIs(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//удаляем старую сессию сохранив ее сначало в переменную

	err = deleteSession(userId)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	in := session.ExpireIn

	if in.Before(time.Now()) {
		http.Error(w, errors.New("Session expired").Error(), http.StatusUnauthorized)
		return
	}

	// сравниваем получены рефреш токен и в токен в сессии

	if refreshToken != session.RefreshToken {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// создаем новые токены
	accessToken, err := jwt.CreateJWTToken(userId)
	newRefreshToken, err := jwt.CreateJWTRefreshToken(userId)

	session.RefreshToken = newRefreshToken
	session.ExpireIn = time.Now().Add(repo.Exp_session)

	// создаем новую сессию с новым токеном
	createSession(*session)

	token := model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	marshal, err := json.Marshal(token)
	if err != nil {
		log.NewLog().Fatal(err)
	}

	fmt.Fprint(w, string(marshal))
}
