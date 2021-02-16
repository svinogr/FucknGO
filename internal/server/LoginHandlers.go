package server

import (
	"FucknGO/internal/jwt"
	"FucknGO/internal/model/user"
	"github.com/gorilla/mux"
	"net/http"
)

// auth user and send jwt token
func auth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email := vars["email"]
	pass := vars["password"]

	if user, err := validUser(email, pass); err != nil {

		jwt.CreateJWT(user.Id)

	}
}

//проверка на валидность юзера в базе
func validUser(email string, pass string) (user.UserModel, error) {
	return user.UserModel{}, nil
}
