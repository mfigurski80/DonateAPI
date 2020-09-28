package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// Start runs server on given addr and port
func Start(addr string) {
	err := store.Init("./data")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	addAuthSubrouter(r)
	addJobSubrouter(r)
	addUserSubrouter(r)

	r.Use(jsonResponseMiddleware)
	r.Use(loggingMiddleware)

	store.Log("Initialized. Serving on http://0.0.0.0" + addr)
	http.ListenAndServe(addr, r)
}
