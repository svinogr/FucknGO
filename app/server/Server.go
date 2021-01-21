package main

import (
	"fmt"
	"log"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная")
}

func main() {
	fileServer := http.FileServer(http.Dir(".app/server/ui/static/"))

	http.HandleFunc("/", mainPage)
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("for start go to:  127.0.0.1:8080/")
}
