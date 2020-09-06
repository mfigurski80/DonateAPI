package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
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

	var user UserRegister
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	users := usersReader.read()
	_, ok := users[user.Username]
	if ok {
		BadRequest(w, "User already exists")
		return
	}

	users[user.Username] = *(newUser(user))
	usersReader.write(users)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

func addAuthSubrouter(r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	authRouter.HandleFunc("/register", handleRegister).Methods(http.MethodPost)
}
