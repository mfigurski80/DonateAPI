package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/state"
)

type postUserStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newUser(u postUserStruct) *state.User {
	return &state.User{
		Username: u.Username,
		Password: state.HashPassword(u.Password),
		Authored: make([]string, 0),
		Running:  make([]string, 0),
	}
}

// POST `/register` creates a new user
func handleRegister(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
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

	var user postUserStruct
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	users := state.UserState.Read()
	_, ok := users[user.Username]
	if ok {
		badRequest(w, "User already exists")
		return
	}

	users[user.Username] = *(newUser(user))
	state.UserState.Write(users)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "createdId": "%s"}`, user.Username)))
}

func addAuthSubrouter(r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	authRouter.HandleFunc("/register", handleRegister).Methods(http.MethodPost)
}
