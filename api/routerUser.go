package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// GET `/{userId}` returns user data (only on self)
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

// PUT `/{userId}` allows updates to user data
func putUser(w http.ResponseWriter, r *http.Request) {
	// auth
	userID := mux.Vars(r)["userId"]
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to update this page")
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
	users[user.Username] = user
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	// respond
	w.Write(makeSuccessResponse("user password has been updated", user.Username))
}

// DELETE `/user` deletes user data and jobs associated with current user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	// auth
	userID := mux.Vars(r)["userId"]
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to delete this page")
		return
	}

	// implement changes
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	delete(users, user.Username)
	for _, ref := range user.Running { // remove references to running
		job := users[ref.User].Authored[ref.Title]
		job.Runner = ""
		users[ref.User].Authored[ref.Title] = job
	}
	for _, job := range user.Authored { // remove authored jobs runners
		if job.Runner == "" {
			continue
		}
		// TODO
	}
	err = store.WriteUsers(users)

	// respond
	w.Write(makeSuccessResponse("user has been deleted", user.Username))
}

func addUserSubrouter(r *mux.Router) {
	userRouter := r.PathPrefix("/{userId}").Subrouter()
	userRouter.HandleFunc("", getUser).Methods(http.MethodGet)
}
