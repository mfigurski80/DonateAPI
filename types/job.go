package types

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
