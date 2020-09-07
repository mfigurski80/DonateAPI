package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GET `/user` returns data on current user
func getUser(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	username, pass, ok := r.BasicAuth()
	if !ok {
		Unauthorized(w)
		return
	}
	user, ok := authUser(username, pass)
	if !ok {
		Unauthorized(w)
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// PUT `/user` allows updates to current user's data
func putUser(w http.ResponseWriter, r *http.Request) {

}

// DELETE `/user` deletes user data and jobs associated with current user
func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func addUserSubrouter(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("", getUser).Methods(http.MethodGet)
	userRouter.HandleFunc("", putUser).Methods(http.MethodPut)
	userRouter.HandleFunc("", deleteUser).Methods(http.MethodDelete)
}
