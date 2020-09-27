package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mfigurski80/DonateAPI/store"
)

type postJobStruct struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func makeJob(j postJobStruct, u store.User) *store.Job {
	return &store.Job{
		Title:         j.Title,
		Description:   j.Description,
		OriginalImage: j.Image,
		Author:        u.Username,
	}
}

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

// POST `/job` creates a new job with given data
func postJob(w http.ResponseWriter, r *http.Request) {
	// auth
	user, ok := authRequest(w, r)
	if !ok {
		return
	}

	// read body
	var postJob postJobStruct
	ok = unmarshalRequestBody(w, r, postJob)
	if !ok {
		return
	}

	// implement changes
	job := makeJob(postJob, user)
	users, err := store.ReadUsers()
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	users[user.Username].Authored[job.Title] = *job
	err = store.WriteUsers(users)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}

	// respond
	w.Write(makeSuccessResponse("job has been created", job.Title))
}

func addJobSubrouter(r *mux.Router) {
	jobRouter := r.PathPrefix("/job").Subrouter()
	jobRouter.HandleFunc("", getJobs).Methods(http.MethodGet)
	jobRouter.HandleFunc("", postJob).Methods(http.MethodPost)
}
