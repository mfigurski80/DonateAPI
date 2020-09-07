package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type PostJobStruct struct {
	Description   string `json:"description"`
	ImageLocation string `json:"imageLocation"`
}

func newJob(s PostJobStruct, author string) *Job {
	return &Job{
		ID:            fmt.Sprintf("%d", time.Now().UnixNano()),
		Description:   s.Description,
		ImageLocation: s.ImageLocation,
		Author:        author,
		Runner:        "",
	}
}

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
	logRequest(r)
}

// POST `/jobs` creates a new job with given data
func postJob(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	username, pass, ok := r.BasicAuth()
	if !ok {
		Unauthorized(w)
		return
	}
	user, ok := authUser(username, pass)
	if !ok {
		Unauthorized(w)
		return
	}

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

	var jobData PostJobStruct
	err = json.Unmarshal(bodyBytes, &jobData)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}
	job := newJob(jobData, user.Username)

	jobs := jobsReader.read()
	_, ok = jobs[job.ID]
	if ok {
		BadRequest(w, "Job already exists")
		return
	}

	jobs[job.ID] = *job
	jobsReader.write(jobs)

	user.Authored = append(user.Authored, job.ID)
	users := usersReader.read()
	users[user.Username] = user
	usersReader.write(users)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "createdId": "%s"}`, job.ID)))
}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()
	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("/{id}", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("", postJob).Methods(http.MethodPost)
}
