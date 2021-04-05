package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/log"
	"html/template"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/mainpage.html")
	conf, err := config.GetConfig()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	base := repo.NewDataBase(conf)
	userRepo := base.User()

	allUser, err := userRepo.FindAllUser()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	files.Execute(w, allUser)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/loginpage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.Execute(w, nil)
}
