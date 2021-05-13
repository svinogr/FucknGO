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
		getAllShopsByAccount(w, r)
	case http.MethodPut:
		updateShop(w, r)
	case http.MethodDelete:
		deleteShopById(w, r)
	}
}

func getAllShopsByAccount(w http.ResponseWriter, r *http.Request) {
	user, err := jwt.GetUserFromContext(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	shopRepo := db.Shop()

	allFinded, err := shopRepo.FindByUserId(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	all := []model.ShopModel{}

	coordRepo := db.Coord()

	for _, el := range *allFinded {
		coordById, err := coordRepo.FindById(el.CoordId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			continue
		}

		m := model.ShopModel{
			Id:       el.Id,
			UserId:   el.UserId,
			CoordLat: strconv.FormatFloat(coordById.CoordLat, 'f', -1, 64),
			CoordLng: strconv.FormatFloat(coordById.CoordLng, 'f', -1, 64),
			Name:     el.Name,
			Address:  el.Address,
		}

		all = append(all, m)
	}

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
	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	coord := repo.CoordModelRepo{}
	coord.CoordLat, err = strconv.ParseFloat(sM.CoordLat, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coord.CoordLng, err = strconv.ParseFloat(sM.CoordLng, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Coord().Create(&coord)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shop := repo.ShopModelRepo{}
	shop.CoordId = coord.Id
	shop.Name = sM.Name
	shop.Address = sM.Address
	shop.UserId = user.Id

	_, err = db.Shop().Create(&shop)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sM.Id = shop.Id

	/*w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)*/
}

//updateShop update shop in db
func updateShop(w http.ResponseWriter, r *http.Request) {
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

	coord := repo.CoordModelRepo{}
	coord.CoordLat, err = strconv.ParseFloat(sM.CoordLat, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coord.CoordLng, err = strconv.ParseFloat(sM.CoordLng, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := repo.NewDataBaseWithConfig()
	defer db.CloseDataBase()

	repoShop := db.Shop()
	shopByid, err := repoShop.FindById(sM.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Id != shopByid.UserId {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repoCoord := db.Coord()
	coordById, err := repoCoord.FindById(shopByid.CoordId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	coordById.CoordLng, err = strconv.ParseFloat(sM.CoordLng, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	coordById.CoordLat, err = strconv.ParseFloat(sM.CoordLat, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shopByid.Address = sM.Address
	shopByid.Name = sM.Name

	_, err = repoShop.Update(shopByid)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*	w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sM)*/
}
