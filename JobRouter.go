package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GET `/jobs` returns list of all *active* jobs (waiting for runners)
func getJobs(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	jobs := jobsReader.read()

	jsonBytes, err := json.Marshal(jobs)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// GET `/jobs/{id}` returns job with given id
func getJob(w http.ResponseWriter, r *http.Request) {

}

// POST `/jobs` creates a new job with given data
func makeJob(w http.ResponseWriter, r *http.Request) {

}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()
	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("/{id}", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("", makeJob).Methods(http.MethodPost)
}
