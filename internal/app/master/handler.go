package master

import (
	"FucknGO/internal/app/config"
	"FucknGO/internal/app/log"
	"FucknGO/internal/app/slave"
	"FucknGO/pkg/requests"
	"FucknGO/pkg/responses"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("home"))
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err.Error(), "data": nil})
		return
	}
}

func listSlaves(w http.ResponseWriter, r *http.Request) {
	data := map[string]responses.Slave{}
	for id, val := range slave.GetSlaves() {
		data[strconv.FormatInt(id, 10)] = responses.Slave{Host: val.Addr}
	}
	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": data})
}

func createSlave(w http.ResponseWriter, r *http.Request) {
	slaveParams := new(requests.SlaveParams)
	if err := render.Bind(r, slaveParams); err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err.Error(), "data": nil})
		return
	}

	conf := config.Get()

	go slave.RunServer(slaveParams)

	msg := fmt.Sprintf("Successfully started slave server on http://%s:%d\n",
		conf.MasterServer.Host, slaveParams.Port)
	logger := log.GetLogger()
	logger.Print(msg)

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": msg})
}

func getSlaveByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	data := map[string]responses.Slave{}
	for id, val := range slave.GetSlaves() {
		data[strconv.FormatInt(id, 10)] = responses.Slave{Host: val.Addr}
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": data[strconv.FormatInt(id, 10)]})
}

func deleteSlave(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURL(r)
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	slaves := slave.GetSlaves()
	err = slaves[id].Close()
	if err != nil {
		render.JSON(w, r, map[string]interface{}{"status": "fail", "error": err, "data": nil})
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok", "error": nil, "data": nil})
}

func getIDFromURL(r *http.Request) (int64, error) {
	if ID := chi.URLParam(r, "ID"); ID != "" {
		id, err := strconv.ParseInt(ID, 10, strconv.IntSize)
		if err != nil {
			return 0, err
		}
		// zero value can be treated as delete all records
		if id == 0 {
			return 0, fmt.Errorf("zero value not allowed")
		}
		// check exist
		slaves := slave.GetSlaves()
		if _, ok := slaves[id]; !ok {
			return 0, fmt.Errorf("slave with id=%d not exist", id)
		}
		return id, err
	}
	return 0, fmt.Errorf("ID parameter not found in url")
}
