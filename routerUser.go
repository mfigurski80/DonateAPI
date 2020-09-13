package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/state"
	"github.com/mfigurski80/DonateAPI/types"
)

// GET `/user` returns data on current user
func getUser(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// authorize
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// write json
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		internalServerError(w, err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// PUT `/user` allows updates to current user's data
func putUser(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// authorize
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// read json
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		unsupportedMediaType(w)
		return
	}
	var newUserData types.PostUserStruct
	err = json.Unmarshal(bodyBytes, &newUserData)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	// update user
	if !(newUserData.Username == user.Username || newUserData.Username == "") {
		badRequest(w, "Username cannot be changed once set")
		return
	}
	if newUserData.Password != "" {
		user.Password = state.HashPassword(newUserData.Password)
	}

	users := state.UserState.Read()
	users[user.Username] = user
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

// DELETE `/user` deletes user data and jobs associated with current user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// authorize
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// delete user
	users := state.UserState.Read()
	delete(users, user.Username)
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

func addUserSubrouter(r *mux.Router) {
	userRouter := r.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("", getUser).Methods(http.MethodGet)
	userRouter.HandleFunc("", putUser).Methods(http.MethodPut)
	userRouter.HandleFunc("", deleteUser).Methods(http.MethodDelete)
}
