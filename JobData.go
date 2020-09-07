package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type Job struct {
	ID            string `json:"id"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	ImageLocation string `json:"imageLocation"`
	Runner        string `json:"runner"`
}

type JobReader struct {
	path  string
	cache map[string]Job
	sync.Mutex
}

func newJobsReader(p string) *JobReader {
	return &JobReader{
		path:  p,
		cache: map[string]Job{},
	}
}

func (r *JobReader) read() map[string]Job {
	r.Lock()
	if len(r.cache) > 0 {
		return r.cache
	}

	file, err := ioutil.ReadFile(r.path)
	if err != nil {
		panic(fmt.Sprintf("Error opening file '%s'", r.path))
	}

	var jobs map[string]Job
	json.Unmarshal(file, &jobs)

	r.cache = jobs
	r.Unlock()
	return jobs
}

func (r *JobReader) write(jobs map[string]Job) {
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

var jobsReader = newJobsReader("data/Jobs.json")
