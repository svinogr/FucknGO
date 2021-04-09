package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Print(r.UserAgent())

	// юзер есть с таким паролем
	db := repo.NewDataBaseWithConfig()
	userRepo := db.User()
	validUser, err := userRepo.GetValidUser(uM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, _ := jwt.CreateJWTToken(validUser.Id)
	refreshToken, _ := jwt.CreateJwtRefreshToken(validUser.Id)

	// есть ли уже сесиия для данного юзера
	sessionRepo := db.Sessions()
	_, err = sessionRepo.GetSessionForUserIdIfIs(validUser.Id)

	if err == nil {
		_, err := sessionRepo.DeleteSessionByUserId(validUser.Id)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}
	// создаем новую сессию
	session := repo.SessionModelRepo{
		UserId:       validUser.Id,
		RefreshToken: refreshToken.Value,
		UserAgent:    r.UserAgent(),
		Fingerprint:  "",
		Ip:           r.RemoteAddr,
		ExpireIn:     time.Now().Add(repo.Exp_session),
		CreatedAt:    time.Now(),
	}

	_, err = sessionRepo.CreateSession(&session)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// добавляем httpOnly Cookie
	jwt.SetCookieWithToken(&w, accessToken)
	jwt.SetCookieWithToken(&w, refreshToken)
}

func updateSession(session *repo.SessionModelRepo) {
	base := repo.NewDataBaseWithConfig()

	sessionRepo := base.Sessions()

	_, err := sessionRepo.UpdateSession(session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func logOut(w http.ResponseWriter, r *http.Request) {
	context := r.Context()
	userId := context.Value(jwt.UserId)

	//TODO проверить есть ли уже сессия ?? рабоатет без проверки

	db := repo.NewDataBaseWithConfig()
	sessionsRepo := db.Sessions()
	_, err := sessionsRepo.DeleteSessionByUserId(uint64(userId.(float64)))

	if err != nil {
		log.NewLog().Fatal(err)
	}

	jwt.DeleteCookie(&w)
}

//  auth/refresh-tokens
func refreshToken(w http.ResponseWriter, r *http.Request) {
	/*tM := model.TokenModel{}
	if err := json.NewDecoder(r.Body).Decode(&tM); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	r.Cookie(model.RefreshTokenName)

	refreshToken := tM.RefreshToken*/

	cookie, err := r.Cookie(model.RefreshTokenName) // получаем куки по ключу рефреш

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	tokenValue := cookie.Value

	claims, err := jwt.GetClaims(tokenValue)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId := uint64(claims[jwt.UserId].(float64))
	expToken := time.Unix(int64(claims[jwt.ExpToken].(float64)), 0)

	// проверяем если срок токена меньше данного момента то
	if expToken.Before(time.Now()) {
		http.Error(w, errors.New("token is expired").Error(), http.StatusUnauthorized)
		return
	}
	base := repo.NewDataBaseWithConfig()
	sessionRepo := base.Sessions()

	sessionOld, err := sessionRepo.GetSessionForUserIdIfIs(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//удаляем старую сессию сохранив ее сначало в переменную

	_, err = sessionRepo.DeleteSessionByUserId(userId)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	in := sessionOld.ExpireIn

	if in.Before(time.Now()) {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// сравниваем получены рефреш токен и в токен в сессии

	if tokenValue != sessionOld.RefreshToken {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// проверяем сессию на клиента и ip

	ok := repo.ValidSession(sessionOld, r)

	if !ok {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// создаем новые токены
	accessToken, err := jwt.CreateJWTToken(userId)
	newRefreshToken, err := jwt.CreateJwtRefreshToken(userId)

	sessionOld.RefreshToken = newRefreshToken.Value
	sessionOld.ExpireIn = time.Now().Add(repo.Exp_session)

	// создаем новую сессию с новым токеном
	sessionRepo.CreateSession(sessionOld)

	jwt.SetCookieWithToken(&w, accessToken)
	jwt.SetCookieWithToken(&w, newRefreshToken)
}
