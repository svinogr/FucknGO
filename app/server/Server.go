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
	fileServer := http.FileServer(http.Dir("./app/server/ui/static")) // ./app/server/ui/static  if start SERVER from root project

	http.HandleFunc("/", mainPage)
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	go startServer()

	 fmt.Println("Successfully started server on http://localhost:8080")

	select {}
}

func startServer() {
	err := http.ListenAndServe(":8080", nil )

	if err != nil {
		log.Fatalf("Error while starting serverErrors: %v", err)
	}
}