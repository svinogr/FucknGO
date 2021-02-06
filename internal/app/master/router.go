package master

import (
	"github.com/go-chi/chi"
	"net/http"
)

func route() http.Handler {
	// create `ServerMux`
	mux := chi.NewRouter()

	//
	mux.Mount("/api", apiRoute())

	return mux
}

func apiRoute() http.Handler {
	r := chi.NewRouter()

	r.Get("/", home)
	r.Mount("/slave", slaveRoute())

	return r
}

func slaveRoute() http.Handler {
	r := chi.NewRouter()

	r.Get("/", listSlaves)
	r.Post("/", createSlave)
	r.Get("/{ID}", getSlaveByID)
	r.Delete("/{ID}", deleteSlave)

	return r
}
