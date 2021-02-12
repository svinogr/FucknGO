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

// newServer creates new servers with input parameters
func NewServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	port := query.Get("port")
	staticPath := query.Get("staticPath")

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

			if data, err := json.Marshal(&s); err != nil {
				fmt.Fprint(w, "not to marshal json")
			} else {
				fmt.Print(data)
				fmt.Fprint(w, string(data))
			}
		}

	} else {
		fmt.Fprint(w, "invalid parameters")
	}

}
