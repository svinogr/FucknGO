package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	user2 "FucknGO/db/user"
	"FucknGO/internal/server/model"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
	case http.MethodDelete:
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var uM = model.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := user2.UserModelRepo{}
	user.Email = uM.Email
	user.Name = uM.Name

	if passwordCrypted, err := bcrypt.GenerateFromPassword([]byte(uM.Password), bcrypt.MinCost); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		user.Password = string(passwordCrypted)
	}

	conf, err := config.GetConfig()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	db := repo.NewDataBase(conf)
	_, err = db.User().CreateUser(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uM.Id = user.Id

	fmt.Fprint(w, user)
}
