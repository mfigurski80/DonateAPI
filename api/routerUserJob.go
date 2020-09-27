package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// GET `/{userId}/{jobId}` returns specific user's job
func getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jobID := vars["jobId"]

	// respond
	job, err := store.ReadJob(userID, jobID)
	if err != nil {
		respondNotFound(w, err.Error())
	}

	bytes, err := json.Marshal(job)
	if err != nil {
		respondInternalServerError(w, err)
	}

	w.Write(bytes)
}

// DELETE `/{userId}/{jobId}` remove specific user's job
func deleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jobID := vars["jobId"]

	// auth
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to update this page")
		return
	}

	// implement changes
	_, ok = user.Authored[jobID]
	if !ok {
		respondNotFound(w, fmt.Sprintf("job '%s' not found on user '%s'", jobID, userID))
	}
	delete(user.Authored, jobID)

	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	users[user.Username] = user
	store.WriteUsers(users)

	// respond
	w.Write(makeSuccessResponse("job has been deleted", jobID))
}

func addUserJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/{jobId}").Subrouter()
	jobRouter.HandleFunc("", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("", deleteJob).Methods(http.MethodDelete)
}
