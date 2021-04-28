package server

import (
	"FucknGO/broker"
	"FucknGO/db/repo"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func userApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodGet:
		getAllUser(w, r)
	case http.MethodDelete:
		deleteUserById(w, r)
	}
}

func deleteUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	userRepo := db.User()

	user := repo.UserModelRepo{Id: id}

	_, err = userRepo.DeleteUser(&user)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getAllUser(w http.ResponseWriter, r *http.Request) {

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
	user.Type = repo.Shop

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
	defer db.CloseDataBase()

	_, err := db.User().CreateUser(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	uM.Id = user.Id

	uM.Password = " "

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
	json.NewEncoder(w).Encode(uM)
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
	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	userRepo := db.User()
	email, err := userRepo.FindUserByEmail(user.Email)

	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			log.NewLog().Fatal(err)
		}
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
