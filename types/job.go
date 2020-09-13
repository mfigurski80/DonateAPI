package types

// Job is a single job datum
type Job struct {
	ID                   string   `json:"id"`
	OriginalImage        string   `json:"originalImage"`
	CompletedImage       string   `json:"completedImage"`
	PartialImages        []string `json:"partialImaged"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	Timestamp            int64    `json:"timestamp"`
	AllowMultipleRunners bool     `json:"allowMultipleRunners"`
	Author               string   `json:"author"`
	Runners              []string `json:"runner"`
}

// ReturnJobStruct is a data type recieved when job is returned to hub
type ReturnJobStruct struct {
	Image       string `json:"image"`
	IsCompleted bool   `json:"isCompleted"`
}

// NewJobStruct is a data type recieved when new job is created
type NewJobStruct struct {
	Title                string `json:"title"`
	Description          string `json:"description"`
	OriginalImage        string `json:"originalImage"`
	AllowMultipleRunners bool   `json:"allowMultipleRunners"`
}
