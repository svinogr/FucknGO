package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

const parseTimeFormat = "2021-05-30"

func stockApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createStock(w, r)
	case http.MethodGet:
		//getAllStockByShop(w, r)
	case http.MethodPut:
		//updateStock(w, r)
	case http.MethodDelete:
		//deleteStockById(w, r)
	}
}

func createStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//TODO возможно стоит перенести проверкю на принадлежность магазина юзеру в мидлвере
	idShop, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sM = model.StockModel{}

	if err = json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	shopRepo := db.Shop()
	shopById, err := shopRepo.FindById(idShop)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shopById.UserId != user.Id {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stock := repo.ShopStockModelRepo{}
	stock.ShopId = idShop
	stock.Title = sM.Title
	stock.Description = sM.Description
	stock.DateStart, err = time.Parse(time.RFC3339, sM.DateStart)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err.Error())
		return
	}

	stock.DateFinish, err = time.Parse(time.RFC3339, sM.DateFinish)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stockRepo := db.ShopStock()
	stockRepo.Create(&stock)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)
}
