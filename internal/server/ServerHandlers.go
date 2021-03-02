package server

import (
	"FucknGO/config"
	"FucknGO/db/repo"
	"FucknGO/internal/server/model"
	"FucknGO/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//main page
func MainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная страница")
}

//test conection
func Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Print("test connect")
	c, err := config.GetConfig()
	if err != nil {
		fmt.Fprint(w, err)
	}

	database := repo.NewDataBase(c)
	err = database.OpenDataBase()

	if err != nil {
		fmt.Fprint(w, err)
	}
	getConfig, err := config.GetConfig()
	usR := repo.NewDataBase(getConfig).User()

	usR.CreateUser(&repo.UserModelRepo{
		Name:     "vasya",
		Password: "123",
		Email:    "123",
	})

	fmt.Fprint(w, "connect")

}

// GetAllServers gets all running servers
func GetAllServers(w http.ResponseWriter, r *http.Request) {
	if fb, err := FabricServer(); err != nil {
		fmt.Fprint(w, err)
	} else {
		servers := fb.servers
		jsonStr := ""

		for _, el := range servers {
			if el == nil {
				continue
			}

			if el.Port() != "" {
				s, err := json.Marshal(model.ServerModel{el.Id(), el.StaticResource(), el.Port(), el.Address(), true})

				if err != nil {
					continue
				}

				jsonStr = jsonStr + string(s)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, jsonStr)
	}
}

// newServer creates new servers
func Server(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateServer(w, r)
	case http.MethodGet:
		GetAllServers(w, r)
	case http.MethodDelete:
		DeleteServerById(w, r)
	}
}

// DeleteServerById deletes by id
func DeleteServerById(w http.ResponseWriter, r *http.Request) {
	var sM model.ServerModel

	vars := mux.Vars(r)

	id, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fb, err := FabricServer()

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	servers := fb.servers

	for _, el := range servers {
		if el == nil {
			continue
		}

		if el.Id() == id {
			err = el.server.Shutdown(context.Background())
			sM.Port = el.Port()
			sM.Address = el.Address()
			sM.StaticResource = el.StaticResource()
			sM.Id = el.Id()
			sM.IsRun = false
			fb.RemoveServer(*el)
			//TODO сделать удаление из спсика серверов
			break
		}
	}

	jsonStr, err := json.Marshal(&sM)

	if err != nil {
		log.NewLog().PrintError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonStr))
}

//Crate new server
func CreateServer(w http.ResponseWriter, r *http.Request) {
	var sM model.ServerModel
	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port := sM.Port
	staticPath := sM.StaticResource

	if port != "" && staticPath != "" {
		fb, err := FabricServer()

		if err != nil {
			log.NewLog().Fatal(err)
		}

		if ser, err := fb.GetNewSlaveServer("0.0.0.0", port, staticPath); err != nil {
			fmt.Fprint(w, err)
		} else {
			go ser.RunServer()

			s := model.ServerModel{ser.Id(), ser.StaticResource(), ser.Port(), ser.Address(), true}

			data, err := json.Marshal(&s)

			if err != nil {
				fmt.Fprint(w, http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(data))
		}

	} else {
		fmt.Fprint(w, http.StatusBadRequest)
	}
}
func Panic(w http.ResponseWriter, r *http.Request) {
	log.NewLog().Fatal(errors.New("panic"))
}
