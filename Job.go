package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Job struct {
	ID            string `json:"id"`
	author        string `json:"author"`
	description   string `json:"description"`
	imageLocation string `json:"imageLocation"`
	runner        string `json:"runner"`
}

type JobReader struct {
	path  string
	cache []Job
	sync.Mutex
}

func newJobsReader(p string) *JobReader {
	return &JobReader{
		path:  p,
		cache: make([]Job, 0),
	}
}

func (r *JobReader) read() []Job {
	r.Lock()
	if len(r.cache) > 0 {
		return r.cache
	}

	file, err := ioutil.ReadFile(r.path)
	if err != nil {
		panic(fmt.Sprintf("Error opening file '%s'", r.path))
	}

	var jobs []Job
	json.Unmarshal(file, &jobs)

	r.cache = jobs
	r.Unlock()
	return jobs
}

func (r *JobReader) write(jobs []Job) {
	file, err := json.MarshalIndent(jobs, "", " ")
	if err != nil {
		panic(fmt.Sprintf("Error Marshaling jobs: %s", jobs))
	}

	err = ioutil.WriteFile(r.path, file, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing to file '%s'", r.path))
	}

	r.Lock()
	r.cache = jobs
	r.Unlock()
}
