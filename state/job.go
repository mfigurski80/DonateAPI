package state

// Job - single job datum
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
