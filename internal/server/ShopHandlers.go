package server

import (
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func shopApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createShop(w, r)
	case http.MethodGet:
		getAllShops(w, r)
	case http.MethodPut:
		updateShop(w, r)
	case http.MethodDelete:
		deleteShopById(w, r)
	}
}

func getAllShops(w http.ResponseWriter, r *http.Request) {
	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	shopRepo := db.Shop()

	all, err := shopRepo.FindAll()

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(all)
}

// deleteShopById deletes shop by id
func deleteShopById(w http.ResponseWriter, r *http.Request) {
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

	shop := repo.ShopModelRepo{Id: id}

	_, err = shopRepo.Delete(&shop)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// createShop creates new shop in db
func createShop(w http.ResponseWriter, r *http.Request) {
	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var sM = model.ShopModel{}
	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shop := repo.ShopModelRepo{}
	shop.CoordId = sM.CoordId
	shop.Name = sM.Name
	shop.Address = sM.Address
	shop.UserId = user.Id

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	_, err = db.Shop().Create(&shop)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sM.Id = shop.Id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)
}

// createShop update shop in db
func updateShop(w http.ResponseWriter, r *http.Request) {
	var sM = model.ShopModel{}
	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shop := repo.ShopModelRepo{}
	shop.Id = sM.Id
	shop.CoordId = sM.CoordId
	shop.Name = sM.Name
	shop.Address = sM.Address

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	_, err := db.Shop().Create(&shop)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)
}
