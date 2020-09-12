package main

import (
	"encoding/json"
	"io/ioutil"
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
	logRequest(r)
	// authorize
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

	// read json
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		BadRequest(w, err.Error())
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		UnsupportedMediaType(w)
		return
	}
	var newUserData PostUserStruct
	err = json.Unmarshal(bodyBytes, &newUserData)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	// add user
	if !(newUserData.Username == user.Username || newUserData.Username == "") {
		BadRequest(w, "Username cannot be changed once set")
		return
	}
	if newUserData.Email != "" {
		user.Email = newUserData.Email
	}
	if newUserData.Password != "" {
		user.Password = hashPassword(newUserData.Password)
	}

	users := usersReader.read()
	users[user.Username] = user
	usersReader.write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

// DELETE `/user` deletes user data and jobs associated with current user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	// authorize
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

	// delete user
	users := usersReader.read()
	delete(users, user.Username)
	usersReader.write(users)

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
