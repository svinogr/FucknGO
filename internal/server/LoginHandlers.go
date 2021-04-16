package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
	"net/http"
)

func auth(w http.ResponseWriter, r *http.Request) {
	var uM model.UserModel

	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// юзер есть с таким паролем
	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	userRepo := db.User()

	user := repo.UserModelRepo{Name: uM.Name, Password: uM.Password, Email: uM.Email}

	validUser, err := userRepo.GetValidUser(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, err := jwt.CreateJWTToken(validUser.Id, model.AccessTokenName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	refreshToken, _ := jwt.CreateJWTToken(validUser.Id, model.RefreshTokenName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// есть ли уже сесиия для данного юзера
	session := repo.SessionModelRepo{UserId: validUser.Id, UserAgent: r.UserAgent(), Ip: r.RemoteAddr}

	_, err = jwt.CreateNewSessionForToken(&session, refreshToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// добавляем httpOnly Cookie
	jwt.SetCookieWithToken(&w, accessToken)
	jwt.SetCookieWithToken(&w, refreshToken)
}

func updateSession(session *repo.SessionModelRepo) {
	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	sessionRepo := db.Sessions()

	_, err := sessionRepo.UpdateSession(session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func logOut(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	user := context.Value(jwt.User).(repo.UserModelRepo)
	//TODO проверить есть ли уже сессия ?? рабоатет без проверки

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	sessionsRepo := db.Sessions()
	_, err := sessionsRepo.DeleteSessionByUserId(user.Id)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	jwt.DeleteCookie(w, model.RefreshTokenName)
	jwt.DeleteCookie(w, model.AccessTokenName)
}

//  auth/refresh-tokens
func refreshToken(w http.ResponseWriter, r *http.Request) {
	sessionByCookieOld, ok := jwt.GetValidSessionByCookie(r) // получаем  сессию по куки

	if !ok {
		http.Error(w, errors.New("session is expired").Error(), http.StatusUnauthorized)
	}

	// создаем новые токены
	accessToken, err := jwt.CreateJWTToken(sessionByCookieOld.UserId, model.AccessTokenName) // не боимся нила так как нил невозможен

	if err != nil {
		log.NewLog().Fatal(err)
	}
	newRefreshToken, err := jwt.CreateJWTToken(sessionByCookieOld.UserId, model.RefreshTokenName)

	if err != nil {
		log.NewLog().Fatal(err)
	}
	// удаляем старую сессию и создалем
	jwt.CreateNewSessionForToken(sessionByCookieOld, newRefreshToken)

	jwt.SetCookieWithToken(&w, accessToken)
	jwt.SetCookieWithToken(&w, newRefreshToken)
}
