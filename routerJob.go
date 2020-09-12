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
	Title                string `json:"title"`
	Description          string `json:"description"`
	OriginalImage        string `json:"originalImage"`
	AllowMultipleRunners bool   `json:"allowMultipleRunners"`
}

func newJob(s postJobStruct, author string) *state.Job {
	time := time.Now().UnixNano()
	return &state.Job{
		ID:                   fmt.Sprintf("%d", time),
		Author:               author,
		Timestamp:            time,
		Title:                s.Title,
		Description:          s.Description,
		OriginalImage:        s.OriginalImage,
		AllowMultipleRunners: s.AllowMultipleRunners,
	}
}

// GET `/jobs` returns list of all *active* jobs (waiting for runners)
func getJobs(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	jobs := state.JobState.Read()

	jsonBytes, err := json.Marshal(jobs)
	if err != nil {
		internalServerError(w, err.Error())
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
		badRequest(w, "Id does not exist")
		return
	}

	// convert to json
	jsonBytes, err := json.Marshal(job)
	if err != nil {
		internalServerError(w, err.Error())
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
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// read json
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		unsupportedMediaType(w)
		return
	}
	var jobData postJobStruct
	err = json.Unmarshal(bodyBytes, &jobData)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	job := newJob(jobData, user.Username)

	// add to jobs
	jobs := state.JobState.Read()
	_, ok = jobs[job.ID]
	if ok {
		badRequest(w, "Job already exists")
		return
	}

	jobs[job.ID] = *job
	state.JobState.Write(jobs)

	user.Authored = append(user.Authored, job.ID)
	users := state.UserState.Read()
	users[user.Username] = user
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "createdId": "%s"}`, job.ID)))
}

// DELETE /jobs/{id}
func deleteJob(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// auth user
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// find referenced job and delete
	id := mux.Vars(r)["id"]
	jobs := state.JobState.Read()
	job, ok := jobs[id]
	if !ok {
		badRequest(w, "Id does not exist")
		return
	}
	if job.Author != user.Username {
		unauthorized(w)
		return
	}
	delete(jobs, job.ID)
	state.JobState.Write(jobs)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

// PUT /jobs/{id}/checkout
func putJobCheckout(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// auth user
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// register runner within job
	id := mux.Vars(r)["id"]
	jobs := state.JobState.Read()
	job, ok := jobs[id]
	if !ok {
		badRequest(w, "Id does not exist")
		return
	}
	if !job.AllowMultipleRunners && len(job.Runners) > 0 {
		badRequest(w, "This job is already being run")
		return
	}
	job.Runners = append(job.Runners, user.Username)
	jobs[id] = job
	state.JobState.Write(jobs)

	// update user ref
	user.Running = append(user.Running, id)
	users := state.UserState.Read()
	users[user.Username] = user
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "checkedId": "%s"}`, id)))
}

// // PUT /jobs/{id}/checkin
// func putJobCheckin(w http.ResponseWriter, r *http.Request) {
// 	state.LogRequest(r)
// 	// auth user
//  user, ok := state.UserState.AuthRequest(r)
// 	if !ok {
// 		unauthorized(w)
// 		return
// 	}

// 	// register runner within job
// 	id := mux.Vars(r)["id"]
// 	jobs := state.JobState.Read()
// 	job, ok := jobs[id]
// 	if !ok {
// 		badRequest(w, "Id does not exist")
// 		return
// 	}
// 	if job.Runner != username {
// 		badRequest(w, "Your are not currently running this job")
// 		return
// 	}
// 	job.Runner = ""
// 	jobs[job.ID] = job
// 	state.JobState.Write(jobs)

// 	// update user ref
// 	for i, jobID := range user.Running {
// 		if jobID != id {
// 			continue
// 		}
// 		user.Running[i] = user.Running[len(user.Running)-1]
// 		user.Running[len(user.Running)-1] = ""
// 		user.Running = user.Running[:len(user.Running)-1]
// 		break
// 	}
// 	users := state.UserState.Read()
// 	users[username] = user
// 	state.UserState.Write(users)

// 	// respond
// 	w.Header().Add("Content-Type", "application/json")
// 	w.Write([]byte(fmt.Sprintf(`{"message": "success", "checkedId": "%s"}`, id)))
// }

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()

	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("", postJob).Methods(http.MethodPost)
	jobRouter.HandleFunc("/{id}", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("/{id}", deleteJob).Methods(http.MethodDelete)
	jobRouter.HandleFunc("/{id}/checkout", putJobCheckout).Methods(http.MethodPost)
	// jobRouter.HandleFunc("/{id}/checkin", putJobCheckin).Methods(http.MethodPut)
}
