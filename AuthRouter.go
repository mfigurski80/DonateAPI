package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newUser(u UserRegister) *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: hashPassword(u.Password),
		Authored: make([]string, 0),
		Running:  make([]string, 0),
	}
}

// POST `/register` creates a new user
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
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "created": "%s"}`, user.Username)))
}

func addAuthSubrouter(r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	authRouter.HandleFunc("/register", handleRegister).Methods(http.MethodPost)
}
