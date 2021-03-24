package server

import (
	"FucknGO/log"
	"html/template"
	"net/http"
)

func webInterface(w http.ResponseWriter, r *http.Request) {
	files, err := template.ParseFiles("ui/web/templates/webinterface.html")

	if err != nil {
		log.NewLog().Fatal(err)
	}

	files.Execute(w, nil)
}
