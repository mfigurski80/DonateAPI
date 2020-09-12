package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski/DonateAPI/state"
)

func main() {

	r := mux.NewRouter()
	addAuthSubrouter(r)
	addUserSubrouter(r)
	addJobSubrouter(r)

	state.Log("Initialized. Serving on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
