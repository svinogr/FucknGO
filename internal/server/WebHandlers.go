package server

import (
	"FucknGO/db/repo"
	"FucknGO/log"
	"fmt"
	"html/template"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	files := template.Must(template.ParseFiles("ui/web/templates/mainpage.html", "ui/web/templates/header.html"))

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	userRepo := db.User()

	allUser, err := userRepo.FindAllUser()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	files.ExecuteTemplate(w, "main", &allUser)
}

func serverPage(w http.ResponseWriter, r *http.Request) {
	files := template.Must(template.ParseFiles("ui/web/templates/serverpage.html", "ui/web/templates/header.html"))

	fabricServer, err := FabricServer()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	servers := fabricServer.servers
	fmt.Print(len(servers))

	files.ExecuteTemplate(w, "server", &servers)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/loginpage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.Execute(w, nil)
}
