package server

import (
	"FucknGO/config"
	"FucknGO/db"
	serModel "FucknGO/internal/model/server"
	"FucknGO/log"
	"encoding/json"
	"fmt"
	"net/http"
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

			s := serModel.ServerModel{ser.Id(), ser.staticResource, ser.port, ser.address}

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
