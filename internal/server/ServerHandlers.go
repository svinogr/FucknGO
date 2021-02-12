package server

import (
	"FucknGO/config"
	"FucknGO/db"
	serModel "FucknGO/internal/model/server"
	"FucknGO/log"
	"context"
	"encoding/json"
	"fmt"
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

	database := db.NewDataBase(c)
	err = database.OpenDataBase()

	if err != nil {
		fmt.Fprint(w, err)
	}

	fmt.Fprint(w, "connect")

}

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
				s, err := json.Marshal(serModel.ServerModel{el.Id(), el.StaticResource(), el.Port(), el.Address()})

				if err != nil {
					continue
				}

				jsonStr = jsonStr + string(s)
			}
		}
		fmt.Fprint(w, jsonStr)
	}
}

//StopServerById stops server by Id
func StopServerById(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if id, err := strconv.Atoi(query.Get("Id")); err != nil {
		fmt.Fprint(w, err)

	} else {
		if fb, err := FabricServer(); err != nil {
			fmt.Fprint(w, err)
		} else {

			servers := fb.servers
			for _, el := range servers {
				if el == nil {
					continue
				}

				if el.Id() == uint64(id) {
					err = el.server.Shutdown(context.Background())
					break
				}
			}
		}
	}
}

// newServer creates new servers
func Server(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateServer(w, r)
	case http.MethodGet:
		GetAllServers(w, r)
	}
}
func CreateServer(w http.ResponseWriter, r *http.Request) {
	var sM serModel.ServerModel
	if err := json.NewDecoder(r.Body).Decode(&sM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port := sM.Port
	staticPath := sM.StaticResource
	fmt.Print(port, staticPath)

	if port != "" && staticPath != "" {
		fb, err := FabricServer()

		if err != nil {
			log.NewLog().Fatal(err)
		}

		if ser, err := fb.GetNewSlaveServer("0.0.0.0", port, staticPath); err != nil {
			fmt.Fprint(w, err)
		} else {
			go ser.RunServer()

			s := serModel.ServerModel{ser.Id(), ser.StaticResource(), ser.Port(), ser.Address()}

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
