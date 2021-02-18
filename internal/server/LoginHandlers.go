package server

import (
	"FucknGO/internal/jwt"
	"FucknGO/internal/model/user"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// auth user and send jwt token

func logPage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("log.html")
	tmp.Execute(w, "done")
}

func auth(w http.ResponseWriter, r *http.Request) {
	var userM user.UserModel
	if err := json.NewDecoder(r.Body).Decode(&userM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := validUser(userM)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	token, _ := jwt.CreateJWT(user.Id)

	c := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(600 * time.Second),
		HttpOnly: true,
	}

	http.SetCookie(w, &c)
	fmt.Fprint(w, "логинься гад")
}

//проверка на валидность юзера в базе
func validUser(user user.UserModel) (user.UserModel, error) {
	//TODO implement with BD

	if user.Password == "1" {
		return user, http.ErrNoCookie
	}

	user.Id = 5
	// end implement
	return user, nil
}
