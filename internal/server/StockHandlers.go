package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/internal/jwt"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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

/*
func validStockByShopByUser(r *http.Request) (idUser uint64, idShop uint64,  idStock uint64, err error)  {

}
*/

func updateStock(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(2024)
	if err != nil {
		fmt.Println(err.Error())
	}

	form := r.MultipartForm
	jsonM := form.Value["json"]

	var sM = model.StockModel{}

	err = json.Unmarshal([]byte(jsonM[0]), &sM) // получем из json обьект

	if err != nil {
		log.NewLog().PrintError(err)
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

	img, header, err := r.FormFile("img") // получаем файл картинку

	if img != nil {
		defer img.Close()
	}

	if err != nil {

		if err.Error() == "http: no such file" { // если файла нет значит картинку надо удалить из базы

			if stock.Img != "-1.jpg" { // если картинка был (по умолчанию она -1) то удаляем ее с диска
				conf, err := config.GetConfig() // конфиг для удаления картинки на диск

				if err != nil {
					log.NewLog().PrintError(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				dst := conf.JsonStr.UiConfig.WWW.StorageImgStock + "/" + stock.Img // получаем адрес картинки для удаления

				err = os.Remove(dst)

				if err != nil {
					log.NewLog().PrintError(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				// ставим картинку по умолчанию
				stock.Img = "-1.jpg"
			}
		} else {
			log.NewLog().PrintError(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var imgName string

	if img != nil { // если есть файл с картинкой
		conf, err := config.GetConfig() // конфиг для записи картинки на диск

		if err != nil {
			log.NewLog().PrintError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// создаем путь куда запсиать картинку
		imgName = strconv.FormatUint(sM.Id, 10) + "." + strings.Split(header.Filename, ".")[1] // имя картинки будет = id stock без разрещения
		//imgName = strconv.FormatUint(sM.Id, 10) // имя картинки будет = id stock
		dst, err := os.OpenFile(conf.JsonStr.UiConfig.WWW.StorageImgStock+"/"+imgName, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			log.NewLog().PrintError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(dst, img)

		if err != nil {
			log.NewLog().PrintError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stock.Img = imgName
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
	sM.Img = stock.Img

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
