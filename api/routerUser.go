package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// GET `/{userId}`
func getUser(w http.ResponseWriter, r *http.Request) {
	// auth
	userID := mux.Vars(r)["userId"]
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to access this page")
		return
	}

	// respond
	bytes, err := json.Marshal(user)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	w.Write(bytes)
}

// PUT `/{userId}`
func putUser(w http.ResponseWriter, r *http.Request) {
	// auth
	userID := mux.Vars(r)["userId"]
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to access this page")
		return
	}

	// read body
	var postUser postUserStruct
	ok = unmarshalRequestBody(w, r, postUser)
	if !ok {
		return
	}

	// implement changes
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	user.Password = store.HashPassword(postUser.Password)
	(*users)[user.Username] = user
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	// respond
	w.Write(makeSuccessResponse("user password has been updated", user.Username))
}

func addUserSubrouter(r *mux.Router) {
	userRouter := r.PathPrefix("/{userId}").Subrouter()
	userRouter.HandleFunc("", getUser).Methods(http.MethodGet)
}
