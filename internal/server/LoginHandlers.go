package server

import (
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
	fmt.Print(r.UserAgent())

	// юзер есть с таким паролем
	validUser, err := getValidUser(uM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	accessToken, _ := jwt.CreateJWTToken(validUser.Id)
	refreshToken, _ := jwt.CreateJWTRefreshToken(validUser.Id)

	// есть ли уже сесиия для данного юзера
	_, err = getSessionForUserIdIfIs(validUser.Id)

	if err == nil {
		err := deleteSession(validUser.Id)

		if err != nil {
			log.NewLog().Fatal(err)
		}
	}
	// создаем новую сессию
	session := repo.SessionModelRepo{
		UserId:       validUser.Id,
		RefreshToken: refreshToken,
		UserAgent:    r.UserAgent(),
		Fingerprint:  "",
		Ip:           r.RemoteAddr,
		ExpireIn:     time.Now().Add(repo.Exp_session),
		CreatedAt:    time.Now(),
	}

	createSession(session)

	// добавляем новые новые токены в ответ
	token := model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func updateSession(session *repo.SessionModelRepo) {
	base := repo.NewDataBaseWithConfig()

	sessionRepo := base.Sessions()

	_, err := sessionRepo.UpdateSession(session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func getSessionForUserIdIfIs(id uint64) (*repo.SessionModelRepo, error) {
	base := repo.NewDataBaseWithConfig()

	sessionRepo := base.Sessions()
	session, err := sessionRepo.FindSessionByUserId(id)

	if err != nil {

		return nil, err
	}

	return session, nil
}

func createSession(session repo.SessionModelRepo) {
	base := repo.NewDataBaseWithConfig()

	sessionRepo := base.Sessions()

	_, err := sessionRepo.CreateSession(&session)

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

// validUser gets valid user by email and password
func getValidUser(user model.UserModel) (*repo.UserModelRepo, error) {
	userRepo := repo.NewDataBaseWithConfig().User()

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
	context := r.Context()
	userId := context.Value(jwt.UserId)

	//TODO проверить есть ли уже сессия ?? рабоатет без проверки
	err := deleteSession(uint64(userId.(float64)))

	if err != nil {
		log.NewLog().Fatal(err)
	}
}

func deleteSession(userId uint64) error {
	base := repo.NewDataBaseWithConfig()

	repoSession := base.Sessions()

	_, err := repoSession.DeleteSessionByUserId(userId)

	return err
}

//  auth/refresh-tokens
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

	sessionOld, err := getSessionForUserIdIfIs(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//удаляем старую сессию сохранив ее сначало в переменную

	err = deleteSession(userId)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	in := sessionOld.ExpireIn

	if in.Before(time.Now()) {
		http.Error(w, errors.New("Session expired").Error(), http.StatusUnauthorized)
		return
	}

	// сравниваем получены рефреш токен и в токен в сессии

	if refreshToken != sessionOld.RefreshToken {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// проверяем сессию на клиента и ip

	ok := validSession(sessionOld, r)

	if !ok {
		http.Error(w, errors.New("session expired").Error(), http.StatusUnauthorized)
		return
	}

	// создаем новые токены
	accessToken, err := jwt.CreateJWTToken(userId)
	newRefreshToken, err := jwt.CreateJWTRefreshToken(userId)

	sessionOld.RefreshToken = newRefreshToken
	sessionOld.ExpireIn = time.Now().Add(repo.Exp_session)

	// создаем новую сессию с новым токеном
	createSession(*sessionOld)

	token := model.TokenModel{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(token)
}

func validSession(session *repo.SessionModelRepo, r *http.Request) bool {
	if session.UserAgent != r.UserAgent() {
		return false
	}
	// данная проверка не работает почему то
	/*	if session.Ip != r.RemoteAddr {
			return false
		}
	*/
	return true
}
