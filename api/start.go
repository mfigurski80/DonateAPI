package api

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// Start runs server on given addr and port
func Start(addr string) {
	file, err := os.OpenFile("./data/users.json", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	file.Close()

	r := mux.NewRouter()
	addAuthSubrouter(r)
	addJobSubrouter(r)
	addUserSubrouter(r)

	r.Use(jsonResponseMiddleware)
	r.Use(loggingMiddleware)

	store.Log("Initialized. Serving on http://0.0.0.0" + addr)
	http.ListenAndServe(addr, r)
}
