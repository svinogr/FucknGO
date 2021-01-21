package main

import (
	"fmt"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task No.2")
}

func main() {
	fmt.Println("for start go to:  127.0.0.0:8080/")

	fileServer := http.FileServer(http.Dir("./app/ui/static/"))

	http.HandleFunc("/", mainPage)
	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":80w80", nil)

	if err != nil {
		fmt.Println(err)
	}
}
