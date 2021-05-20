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
		shopAccountPage(w, user)
	case repo.Shop:
		shopAccountPage(w, user)
	case repo.Client:
		clientPage(user)
	}
}

func clientPage(user repo.UserModelRepo) {

}

func shopAccountPage(w http.ResponseWriter, user repo.UserModelRepo) {
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
	files, err := template.ParseFiles("ui/web/templates/newshoppage.html", "ui/web/templates/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.ExecuteTemplate(w, "newshoppage", nil)
}

func newStockPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idShop, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stock := model.StockModel{}
	stock.ShopId = idShop

	files, err := template.ParseFiles("ui/web/templates/newstockpage.html", "ui/web/templates/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.ExecuteTemplate(w, "newstockpage", stock)
}

func shopPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	shopRepo := db.Shop()
	shopById, err := shopRepo.FindById(id)

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if shopById.UserId != user.Id {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stockRepo := db.ShopStock()
	stoks, err := stockRepo.FindByShop(shopById)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sTArray := []model.StockModel{}

	for _, el := range *stoks {
		sT := model.StockModel{}
		sT.Id = el.Id
		sT.ShopId = el.ShopId
		sT.Title = el.Title
		sT.Description = el.Description
		sT.DateStart = el.DateStart.String()
		sT.DateFinish = el.DateStart.String()

		sTArray = append(sTArray, sT)
	}

	coordRepo := db.Coord()
	coordById, err := coordRepo.FindById(shopById.CoordId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sM := model.ShopModel{}
	sM.Id = shopById.Id
	sM.UserId = shopById.UserId
	sM.CoordLng = strconv.FormatFloat(coordById.CoordLng, 'f', -1, 64)
	sM.CoordLat = strconv.FormatFloat(coordById.CoordLat, 'f', -1, 64)
	sM.Address = shopById.Address
	sM.Name = shopById.Name
	sM.Stocks = sTArray
	sM.Id = shopById.Id

	files, err := template.ParseFiles("ui/web/templates/shoppage.html", "ui/web/templates/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	files.ExecuteTemplate(w, "shoppage", &sM)
}

func updateShopPage(w http.ResponseWriter, r *http.Request) {
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
	defer db.CloseDataBase()
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

	files, err := template.ParseFiles("ui/web/templates/changeshoppage.html", "ui/web/templates/header.html")

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

	files.ExecuteTemplate(w, "changeshoppage", sM)
}

func stockPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idShop, err := strconv.ParseUint(vars["id_shop"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	idStock, err := strconv.ParseUint(vars["id_stock"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := repo.NewDataBaseWithConfig()
	stockRepo := db.ShopStock()
	stock := repo.ShopStockModelRepo{Id: idStock}
	_, err = stockRepo.FindById(&stock)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if stock.ShopId != idShop {
		log.NewLog().PrintError(err)
		http.Error(w, errors.New("user denied  to this action ").Error(), http.StatusBadRequest)
		return
	}

	shopRepo := db.Shop()
	shopById, err := shopRepo.FindById(idShop)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Id != shopById.UserId {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files, err := template.ParseFiles("ui/web/templates/stockpage.html", "ui/web/templates/header.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	layout := "2006-01-02"
	sM := model.StockModel{}
	sM.ShopId = stock.ShopId
	sM.Id = stock.Id
	sM.Title = stock.Title
	sM.Description = stock.Description
	sM.DateStart = stock.DateStart.Format(layout)
	sM.DateFinish = stock.DateFinish.Format(layout)
	sM.Img = stock.Img

	files.ExecuteTemplate(w, "stockpage", sM)
}
