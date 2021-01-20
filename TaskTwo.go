package main

import (
	"fmt"
	"net/http"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task No.2")
}

func main() {
	http.HandleFunc("/", mainPage)
	http.ListenAndServe(":8080", nil)
}
