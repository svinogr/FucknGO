package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная")
}

func main() {
	http.HandleFunc("/", mainPage)

	setupStaticResource()

	go startServer()

	fmt.Println("Successfully started server on http://localhost:8080")

	select {}
}

func setupStaticResource() {
	path := "./app/server/ui/static"

	_, err := os.Stat(path)

	if err != nil {
		os.MkdirAll(path, 0777)
	}

	fileServer := http.FileServer(http.Dir(path))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))
}

func startServer() {
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error while starting serverErrors: %v", err)
	}
}
