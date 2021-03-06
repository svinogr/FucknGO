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

// createUser creates new user in db
func createUser(w http.ResponseWriter, r *http.Request) {
	var uM = model.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := repo.UserModelRepo{}
	user.Email = uM.Email
	user.Name = uM.Name
	user.Password = uM.Password

	if !validRegInfo(user) {
		http.Error(w, errors.New("Not valid register data").Error(), http.StatusBadRequest)
		return
	}

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

	uM.Password = " "

	createJWT, err := jwt.CreateJWT(uM.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uM.Token = createJWT
	jsonStr, err := json.Marshal(&uM)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonStr))
}

// validRegInfo validation register data
func validRegInfo(user repo.UserModelRepo) bool {
	if user.Password == "" {
		return false
	}

	if user.Email == "" {
		return false
	}

	if user.Name == "" {
		return false
	}

	config, err := config.GetConfig()

	if err != nil {
		return false
	}

	userRepo := repo.NewDataBase(config).User()

	email, err := userRepo.FindUserByEmail(user.Email)

	if email != nil {
		return false
	}

	name, err := userRepo.FindUserByName(user.Name)

	if name != nil {
		return false
	}

	return true
}
