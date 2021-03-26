package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/log"
	"fmt"
	"html/template"
	"net/http"
)

func webInterface(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/webinterface.html")
	conf, err := config.GetConfig()

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	base := repo.NewDataBase(conf)
	userRepo := base.User()

	allUser, err := userRepo.FindAllUser()

	if err != nil {
		log.NewLog().Fatal(err)
	}
	fmt.Print(allUser)
	files.Execute(w, allUser)
}
