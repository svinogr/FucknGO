package slave

import (
	"net/http"
)

func route(dir string) http.Handler {
	// create `ServerMux`
	mux := http.NewServeMux()

	// create a default route handler
	mux.HandleFunc("/", home)

	// create a default route handler
	mux.Handle("/static/", makeStaticHandler(dir))

	return mux
}
