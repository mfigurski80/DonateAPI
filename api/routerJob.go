package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

// GET `/job` returns list of all *active* jobs (waiting for runners)
func getJobs(w http.ResponseWriter, r *http.Request) {
	// filter data
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	var jobs = make([]store.Job, 0)
	for _, user := range users {
		for _, job := range user.Authored {
			if job.Runner == "" {
				jobs = append(jobs, job)
			}
		}
		if len(jobs) > 100 {
			break
		}
	}

	// respond
	bytes, err := json.Marshal(jobs)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	w.Write(bytes)
}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()
	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
}
