package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/state"
	"github.com/mfigurski80/DonateAPI/types"
)

func newJob(s types.NewJobStruct, author string) *types.Job {
	time := time.Now().UnixNano()
	return &types.Job{
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
	var jobData types.NewJobStruct
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

	// remove all runner and author references
	users := state.UserState.Read()
	for _, runner := range job.Runners {
		u := users[runner]
		u.Running = remove(u.Running, find(u.Running, job.ID))
		users[runner] = u
	}
	author := users[job.Author]
	author.Authored = remove(author.Authored, find(author.Authored, job.ID))
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

// PUT /jobs/{id}/take
func putJobTake(w http.ResponseWriter, r *http.Request) {
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

// PUT /jobs/{id}/return
// checks back in the sent image and disassociates user from it
func putJobReturn(w http.ResponseWriter, r *http.Request) {
	state.LogRequest(r)
	// auth user
	user, ok := state.UserState.AuthRequest(r)
	if !ok {
		unauthorized(w)
		return
	}

	// read return image data
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
	var returnData types.ReturnJobStruct
	err = json.Unmarshal(bodyBytes, &returnData)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	// find referenced job and user runner reference
	id := mux.Vars(r)["id"]
	jobs := state.JobState.Read()
	job, ok := jobs[id]
	if !ok {
		badRequest(w, "Job with this id does not exist")
		return
	}
	runnerIndex := find(job.Runners, user.Username)
	if runnerIndex < 0 {
		badRequest(w, "User is not listed as a runner for this resource")
		return
	}

	// append return image data to current job
	if job.CompletedImage != "" {
		badRequest(w, "This job has already been completed")
		return
	}
	if returnData.IsCompleted {
		job.CompletedImage = returnData.Image
	} else {
		job.PartialImages = append(job.PartialImages, returnData.Image)
	}

	// remove referenced job and user runner reference
	job.Runners = remove(job.Runners, runnerIndex)
	jobs[job.ID] = job
	state.JobState.Write(jobs)
	user.Running = remove(user.Running, find(user.Running, job.ID))
	users := state.UserState.Read()
	users[user.Username] = user
	state.UserState.Write(users)

	// respond
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "success", "checkedId": "%s"}`, job.ID)))
}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/jobs").Subrouter()

	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("", postJob).Methods(http.MethodPost)
	jobRouter.HandleFunc("/{id}", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("/{id}", deleteJob).Methods(http.MethodDelete)
	jobRouter.HandleFunc("/{id}/take", putJobTake).Methods(http.MethodPost)
	jobRouter.HandleFunc("/{id}/return", putJobReturn).Methods(http.MethodPut)
}
