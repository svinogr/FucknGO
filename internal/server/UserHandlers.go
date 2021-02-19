package server

import (
	"FucknGO/internal/server/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func user(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
	case http.MethodDelete:
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var uM = model.UserModel{}
	if err := json.NewDecoder(r.Body).Decode(&uM); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, uM)
}
