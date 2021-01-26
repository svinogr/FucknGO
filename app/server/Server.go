package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Controller function for main page
func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Главная")
}

type Server struct {
	port int8
}

type HandleFunction struct {
	Url string
	Function
}

func (s *Server) registreHandlerFunction(function HandleFunction) error {
	http.Handle(function.Url, function.Function)

}

func mainNo() {
	http.HandleFunc("/", mainPage)

	err := setupStaticResource()

	if err != nil {
		log.Fatal("it doesnt make dir for static resources")
	}

	go startServer()

	fmt.Println("Successfully started server on http://localhost:8080")

	select {}
}

// setupStaticResource creates a directory named path: "./app/server/ui/static",
// along with any necessary parents, and returns nil,
// or else returns an error.
// If path is already a directory, setupStaticResource does nothing
// and returns nil.
func setupStaticResource() error {
	path := "./app/server/ui/static"

	_, err := os.Stat(path)

	if err != nil {
		err = os.MkdirAll(path, 0777)

		if err != nil {
			return err
		}
	}

	fileServer := http.FileServer(http.Dir(path))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	return nil
}

// startServer starts server on port 8080
func startServer() {
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("Error while starting serverErrors: %v", err)
	}
}
