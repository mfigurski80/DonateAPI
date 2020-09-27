package api

import (
	"encoding/json"
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

func addUserJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/{jobId}").Subrouter()
	jobRouter.HandleFunc("", getJob).Methods(http.MethodGet)
}
