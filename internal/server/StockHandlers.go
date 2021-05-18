package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
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
		updateStock(w, r)
	case http.MethodDelete:
		deleteStockById(w, r)
	}
}

func updateStock(w http.ResponseWriter, r *http.Request) {
	var sM = model.StockModel{}

	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, errors.New("user denied  to this action ").Error(), http.StatusBadRequest)
		return
	}

	stock := repo.ShopStockModelRepo{
		Id: idStock,
	}

	stockRepo := db.ShopStock()
	_, err = stockRepo.FindById(&stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shopById.Id != stock.ShopId {
		http.Error(w, errors.New("stock denied  to this shop ").Error(), http.StatusBadRequest)
		return
	}

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

	_, err = stockRepo.Update(&stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sM.Id = stock.Id
	sM.ShopId = stock.ShopId

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)

}

func deleteStockById(w http.ResponseWriter, r *http.Request) {
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

	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, errors.New("user denied  to this action ").Error(), http.StatusBadRequest)
		return
	}

	stock := repo.ShopStockModelRepo{
		Id: idStock,
	}

	stockRepo := db.ShopStock()
	_, err = stockRepo.FindById(&stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shopById.Id != stock.ShopId {
		http.Error(w, errors.New("stock denied  to this shop ").Error(), http.StatusBadRequest)
		return
	}

	stockRepo.Delete(&stock)
}

func HaveUserThisShop(r *http.Request) (bool, error) {
	/*	vars := mux.Vars(r)
		//TODO как то надо в этом методе инкапсулировать говно из остальных
		idShop, err := strconv.ParseUint(vars["id_shop"], 10, 32)

		if err != nil {
			log.NewLog().PrintError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := jwt.GetUserFromContext(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}*/
	return false, nil
}

func createStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//TODO возможно стоит перенести проверкю на принадлежность магазина юзеру в мидлвере
	idShop, err := strconv.ParseUint(vars["id_shop"], 10, 32)

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

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	shopRepo := db.Shop()
	shopById, err := shopRepo.FindById(idShop)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shopById.UserId != user.Id {
		http.Error(w, errors.New("user denied  to this action ").Error(), http.StatusBadRequest)
		return
	}

	var sM = model.StockModel{}

	if err = json.NewDecoder(r.Body).Decode(&sM); err != nil {
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
	_, err = stockRepo.Create(&stock)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)
}
