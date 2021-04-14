package server

import (
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

//test conection
/*func Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Print("test connect")

	database := repo.NewDataBaseWithConfig()
	err := database.OpenDataBase()

	if err != nil {
		fmt.Fprint(w, err)
	}

	usR := repo.NewDataBaseWithConfig().User()

	usR.CreateUser(&repo.UserModelRepo{
		Name:     "vasya",
		Password: "123",
		Email:    "123",
	})

	fmt.Fprint(w, "connect")

}
*/
// getAllServers gets all running servers
func getAllServers(w http.ResponseWriter, r *http.Request) {
	if fb, err := FabricServer(); err != nil {
		fmt.Fprint(w, err)
	} else {
		servers := fb.servers

		serverToJSon := []model.ServerModel{}
		for _, el := range servers {
			if el == nil {
				continue
			}

			server := model.ServerModel{
				Id:             el.Id(),
				StaticResource: el.StaticResource(),
				Port:           el.Port(),
				Address:        el.Address(),
				IsRun:          true,
			}
			serverToJSon = append(serverToJSon, server)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serverToJSon)
	}
}

// newServer creates new servers
func serverApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createServer(w, r)
	case http.MethodGet:
		getAllServers(w, r)
	case http.MethodDelete:
		deleteServerById(w, r)
	}
}

// deleteServerById deletes by id
func deleteServerById(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sM)
}

// createServer creates new serverApi
func createServer(w http.ResponseWriter, r *http.Request) {
	var sM model.ServerModel
	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверка что порт нормальный
	if _, err := strconv.Atoi(sM.Port); err != nil {
		http.Error(w, "port is bad", http.StatusBadRequest)
		return
	}
	// TODO проверка адреса что он без слешей

	fb, err := FabricServer()

	if err != nil {
		log.NewLog().Fatal(err)
	}

	slaveServer, err := fb.GetNewSlaveServer("0.0.0.0", sM.Port)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go log.NewLog().PrintCommon(slaveServer.RunServer().Error())

	sM.Id = slaveServer.id
	sM.Address = slaveServer.address
	sM.IsRun = true

	s := model.ServerModel{slaveServer.Id(), slaveServer.StaticResource(), slaveServer.Port(), slaveServer.Address(), true}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

//test function for  panic handler
func Panic(w http.ResponseWriter, r *http.Request) {
	log.NewLog().Fatal(errors.New("panic"))
}
