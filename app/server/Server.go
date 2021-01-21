package main

import (
	"fmt"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task No.2")
}

func main() {
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	http.HandleFunc("/", mainPage)
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.ListenAndServe(":8080", nil)
}
