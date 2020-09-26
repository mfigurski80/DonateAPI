package store

import "fmt"

// Job is a single job datum
type Job struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	OriginalImage  string `json:"originalImage"`
	CompletedImage string `json:"completedImage"`
	Timestamp      int64  `json:"timestamp"`
	Author         string `json:"author"`
	Runner         string `json:"runner"`
}

// JobMap maps job title to job data
type JobMap map[string]Job

var jobPath = "./data/jobs.json"

// ReadJob reads and returns job from given username of given title
func ReadJob(username string, title string) (Job, error) {
	users, err := ReadUsers()
	if err != nil {
		return Job{}, err
	}

	user, ok := (*users)[username]
	if !ok {
		return Job{}, fmt.Errorf("Username '%s' not found", username)
	}

	job, ok := user.Authored[title]
	if !ok {
		return Job{}, fmt.Errorf("Job '%s' by user '%s' not found", title, username)
	}

	return job, nil
}

// WriteJob writes given job to proper user in store file
func WriteJob(job Job) error {
	users, err := ReadUsers()
	if err != nil {
		return err
	}

	user, ok := (*users)[job.Author]
	if !ok {
		return fmt.Errorf("Username '%s' not found", job.Author)
	}

	user.Authored[job.Title] = job
	(*users)[job.Author] = user
	WriteUsers(users)

	return nil
}
