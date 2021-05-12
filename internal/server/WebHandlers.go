package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
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

func changeShopPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	db := repo.NewDataBaseWithConfig()
	shopRepo := db.Shop()
	shopById, err := shopRepo.FindById(id)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shopById.UserId != user.Id {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coordRepo := db.Coord()
	coordById, err := coordRepo.FindById(shopById.CoordId)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files, err := template.ParseFiles("ui/web/templates/changeshoppage.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	sM := model.ShopModel{
		Id:       id,
		UserId:   user.Id,
		CoordLat: strconv.FormatFloat(coordById.CoordLat, 'f', -1, 64),
		CoordLng: strconv.FormatFloat(coordById.CoordLng, 'f', -1, 64),
		Name:     shopById.Name,
		Address:  shopById.Address,
	}

	files.Execute(w, sM)
}
