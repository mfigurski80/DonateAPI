package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

type returnJobStruct struct {
	Image string `json:"image"`
}

func makeReturnedJob(data returnJobStruct, job *store.Job) *store.Job {
	if data.Image != "" {
		(*job).CompletedImage = data.Image
	}
	(*job).Runner = ""
	return job
}

// GET `/{userId}/{jobId}` returns specific user's job
func getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jobID := vars["jobId"]

	// respond
	job, err := store.ReadJob(userID, jobID)
	if err != nil {
		respondNotFound(w, err.Error())
		return
	}

	bytes, err := json.Marshal(job)
	if err != nil {
		respondInternalServerError(w, err)
		return
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
	if _, ok = user.Authored[jobID]; !ok {
		respondNotFound(w, fmt.Sprintf("job '%s' not found on user '%s'", jobID, userID))
		return
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

// PUT `/{userId}/{jobId}/take` marks job as being run
func takeJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jobID := vars["jobId"]

	// auth
	rUser, ok := authRequest(w, r)
	if !ok {
		return
	}

	// implement changes
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	aUser, ok := users[userID]
	if !ok {
		respondNotFound(w, fmt.Sprintf("user '%s' not found", userID))
		return
	}
	job, ok := aUser.Authored[jobID]
	if !ok {
		respondNotFound(w, fmt.Sprintf("job '%s' not found on user '%s'", jobID, userID))
		return
	}
	if job.Runner != "" {
		respondBadRequest(w, fmt.Sprintf("job '%s' is already being run by user '%s'", jobID, job.Runner))
		return
	}
	if job.CompletedImage != "" {
		respondBadRequest(w, fmt.Sprintf("job '%s' is already completed", jobID))
		return
	}
	job.Runner = rUser.Username
	aUser.Authored[jobID] = job
	rUser.Running = append(rUser.Running, store.JobReference{User: job.Author, Title: job.Title})

	users[aUser.Username] = aUser
	users[rUser.Username] = rUser
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	// respond
	w.Write(makeSuccessResponse("job taken", jobID))
}

// PUT `/{userId}/{jobId}/return` marks job as not being run and maybe completed
func returnJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	jobID := vars["jobId"]

	// auth
	rUser, ok := authRequest(w, r)
	if !ok {
		return
	}

	// read data
	var data returnJobStruct
	ok = unmarshalRequestBody(w, r, &data)
	if !ok {
		return
	}

	// implementChanges
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	aUser, ok := users[userID]
	if !ok {
		respondNotFound(w, fmt.Sprintf("user '%s' not found", userID))
		return
	}
	job, ok := aUser.Authored[jobID]
	if !ok {
		respondNotFound(w, fmt.Sprintf("job '%s' not found on user '%s'", jobID, userID))
		return
	}
	if job.Runner != rUser.Username {
		respondUnauthorized(w, fmt.Sprintf("user '%s' is not currently running this job", rUser.Username))
		return
	}
	makeReturnedJob(data, &job)
	aUser.Authored[jobID] = job
	rUser.Running = removeJobReferenceListValue(rUser.Running, store.JobReference{User: job.Author, Title: job.Title})
	users[aUser.Username] = aUser
	users[rUser.Username] = rUser
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	// respond
	w.Write(makeSuccessResponse("job returned", jobID))
}

func addUserJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/{jobId}").Subrouter()
	jobRouter.HandleFunc("", getJob).Methods(http.MethodGet)
	jobRouter.HandleFunc("", deleteJob).Methods(http.MethodDelete)
	jobRouter.HandleFunc("/take", takeJob).Methods(http.MethodPut)
	jobRouter.HandleFunc("/return", returnJob).Methods(http.MethodPut)
}
