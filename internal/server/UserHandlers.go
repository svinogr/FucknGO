package server

import (
	"FucknGO/broker"
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

	db := repo.NewDataBaseWithConfig()
	_, err := db.User().CreateUser(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uM.Id = user.Id

	uM.Password = " "

	createJWT, err := jwt.CreateJWTToken(uM.Id)

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
	// постановка сообщения в очередь для отправки юзеру на почту
	err = broker.PublishMessage(broker.MailMessage{
		Name:     uM.Name,
		Email:    uM.Email,
		Password: uM.Password,
	})

	if err != nil {
		log.NewLog().PrintCommon(err.Error())
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

	userRepo := repo.NewDataBaseWithConfig().User()

	email, err := userRepo.FindUserByEmail(user.Email)

	if err != nil {
		log.NewLog().Fatal(err)
	}

	if email != nil {
		return false
	}

	name, err := userRepo.FindUserByName(user.Name)

	if name != nil {
		return false
	}

	return true
}
