package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {

}

func addSubrouter(r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	authRouter.HandleFunc("/login", handleLogin).Methods(http.MethodPost)
}
