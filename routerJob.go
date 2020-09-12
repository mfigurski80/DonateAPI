package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/state"
)

type postJobStruct struct {
	Description   string `json:"description"`
	ImageLocation string `json:"imageLocation"`
}

func newJob(s PostJobStruct, author string) *state.Job {
	return &state.Job{
		ID:            fmt.Sprintf("%d", time.Now().UnixNano()),
		Description:   s.Description,
		ImageLocation: s.ImageLocation,
		Author:        author,
		Runner:        "",
	}
}

// GET `/jobs` returns list of all *active* jobs (waiting for runners)
func getJobs(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	jobs := state.JobState.Read()

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
	state.LogRequest(r)

	// find referenced job
	id := mux.Vars(r)["id"]
	jobs := state.JobState.Read()
	job, ok := jobs[id]
	if !ok {
		BadRequest(w, "Id does not exist")
		return
	}

	// convert to json
	jsonBytes, err := json.Marshal(job)
	if err != nil {
		InternalServerError(w, err.Error())
		return
	}

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// POST `/jobs` creates a new job with given data
func postJob(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// auth user
	username, pass, ok := r.BasicAuth()
	if !ok {
		Unauthorized(w)
		return
	}
	user, ok := state.UserState.AuthUser(username, pass)
	if !ok {
		Unauthorized(w)
		return
	}

	// read json
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

	// add to jobs
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

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "createdId": "%s"}`, job.ID)))
}

// DELETE /jobs/{id}
func deleteJob(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	// auth user
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

	// find referenced job and delete
	id := mux.Vars(r)["id"]
	jobs := jobsReader.read()
	job, ok := jobs[id]
	if !ok {
		BadRequest(w, "Id does not exist")
		return
	}
	if job.Author != user.Username {
		Unauthorized(w)
		return
	}
	delete(jobs, job.ID)
	jobsReader.write(jobs)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

// PUT /jobs/{id}/checkout
func putJobCheckout(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	// auth user
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

	// register runner within job
	id := mux.Vars(r)["id"]
	jobs := jobsReader.read()
	job, ok := jobs[id]
	if !ok {
		BadRequest(w, "Id does not exist")
		return
	}
	if job.Runner != "" {
		BadRequest(w, "This job is already being run")
		return
	}
	job.Runner = username
	jobs[id] = job
	jobsReader.write(jobs)

	// update user ref
	user.Running = append(user.Running, id)
	users := usersReader.read()
	users[username] = user
	usersReader.write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "checkedId": "%s"}`, id)))
}

// PUT /jobs/{id}/checkin
func putJobCheckin(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	// auth user
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

	// register runner within job
	id := mux.Vars(r)["id"]
	jobs := jobsReader.read()
	job, ok := jobs[id]
	if !ok {
		BadRequest(w, "Id does not exist")
		return
	}
	if job.Runner != username {
		BadRequest(w, "Your are not currently running this job")
		return
	}
	job.Runner = ""
	jobs[job.ID] = job
	jobsReader.write(jobs)

	// update user ref
	for i, jobId := range user.Running {
		if jobId != id {
			continue
		}
		user.Running[i] = user.Running[len(user.Running)-1]
		user.Running[len(user.Running)-1] = ""
		user.Running = user.Running[:len(user.Running)-1]
		break
	}
	users := usersReader.read()
	users[username] = user
	usersReader.write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "checkedId": "%s"}`, id)))
}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()

	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("", postJob).Methods(http.MethodPost)
	jobRouter.HandleFunc("/{id}", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("/{id}", deleteJob).Methods(http.MethodDelete)
	jobRouter.HandleFunc("/{id}/checkin", putJobCheckin).Methods(http.MethodPut)
	jobRouter.HandleFunc("/{id}/checkout", putJobCheckout).Methods(http.MethodPost)
}
