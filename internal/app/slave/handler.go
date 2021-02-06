package slave

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, "hello\n")
	} else {
		http.NotFound(w, r)
		return
	}
}

func makeStaticHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	return http.StripPrefix("/static/", fileServer)
}
