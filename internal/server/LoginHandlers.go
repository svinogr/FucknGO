package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

// auth user and send jwt token
// test handler for aut from html form
func logPage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("log.html")
	tmp.Execute(w, "done")
}

// auth responses with token if log is success
func auth(w http.ResponseWriter, r *http.Request) {
	var uM model.UserModel

	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validUser, err := getValidUser(uM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, _ := jwt.CreateJWT(validUser.Id)

	/*	c := http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(600 * time.Second),
			HttpOnly: true,
		}

		http.SetCookie(w, &c)*/

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{ \"token\": \""+token+"\"}")
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
	value := r.Context().Value(jwt.USER_ID)

	if value == nil {
		return nil, errors.New("Not found id")
	}

	return value, nil
}
