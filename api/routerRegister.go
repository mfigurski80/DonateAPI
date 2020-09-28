package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

type postUserStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func makeUser(u postUserStruct) *store.User {
	return &store.User{
		Username: u.Username,
		Password: store.HashPassword(u.Password),
		Authored: make(store.JobMap, 0),
		Running:  []store.JobReference{},
	}
}

// POST `/register` creates a new user
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var postUser postUserStruct
	ok := unmarshalRequestBody(w, r, &postUser)
	if !ok { // response is handled
		return
	}

	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	_, ok = users[postUser.Username]
	if ok {
		respondUnauthorized(w, fmt.Sprintf("user %s already exists", postUser.Username))
		return
	}
	users[postUser.Username] = *makeUser(postUser)
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	w.Write(makeSuccessResponse("user has been created", postUser.Username))
}

func addAuthSubrouter(r *mux.Router) {
	authRouter := r.NewRoute().Subrouter()
	authRouter.HandleFunc("/register", handleRegister).Methods(http.MethodPost)
}
