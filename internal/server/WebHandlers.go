package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/log"
	"errors"
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
	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Type != repo.Admin {
		http.Error(w, errors.New("access denied").Error(), http.StatusForbidden)
		return
	}

	files := template.Must(template.ParseFiles("ui/web/templates/serverpage.html", "ui/web/templates/header.html"))

	fabricServer, err := FabricServer()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	servers := fabricServer.servers

	err = files.ExecuteTemplate(w, "server", &servers)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/loginpage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.Execute(w, nil)
}

func newuser(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/newuserpage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.Execute(w, nil)
}

func accountPage(w http.ResponseWriter, r *http.Request) {
	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch user.Type {
	case repo.Admin:
		shopPage(w, user)
	case repo.Shop:
		shopPage(w, user)
	case repo.Client:
		clientPage(user)
	}
}

func clientPage(user repo.UserModelRepo) {

}

func shopPage(w http.ResponseWriter, user repo.UserModelRepo) {
	files := template.Must(template.ParseFiles("ui/web/templates/shopaccountpage.html", "ui/web/templates/header.html"))

	db := repo.NewDataBaseWithConfig()

	shops, err := db.Shop().FindByUserId(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = files.ExecuteTemplate(w, "shops", &shops)
}

func newShopPage(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/newshoppage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.Execute(w, nil)
}
