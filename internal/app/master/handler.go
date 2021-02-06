package master

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/db"
	"FucknGO/internal/app/slave"
	"FucknGO/pkg/requests"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("home"))
}

func listSlaves(w http.ResponseWriter, r *http.Request) {
	slaves, err := db.GetList()
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": slaves})
}

func createSlave(w http.ResponseWriter, r *http.Request) {
	slaveParams := new(requests.SlaveParams)
	if err := render.Bind(r, slaveParams); err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err.Error(), "data": nil})
		return
	}

	conf := config.Get()

	go slave.RunServer(slaveParams)
	log.Printf("Successfully started slave server on http://%s:%d\n",
		conf.MasterServer.Host, slaveParams.Port)

	slaveRecord := db.Slave{
		Port:      slaveParams.Port,
		StaticDir: slaveParams.StaticDir,
	}

	if err := db.Create(&slaveRecord); err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err.Error(), "data": nil})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": slaveRecord})
}

func getSlaveByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	slaveRecord, err := db.GetByID(id)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": slaveRecord})
}

func deleteSlave(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	err = db.Delete(id)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": nil})
}

func getIDFromURL(r *http.Request) (uint, error) {
	if ID := chi.URLParam(r, "ID"); ID != "" {
		id, err := strconv.ParseUint(ID, 10, strconv.IntSize)
		if err != nil {
			return 0, err
		}
		// zero value can be treated as delete all records
		if id == 0 {
			return 0, fmt.Errorf("zero value not allowed")
		}
		return uint(id), err
	}
	return 0, fmt.Errorf("ID parameter not found in url")
}
