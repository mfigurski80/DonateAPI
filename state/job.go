package state

// Job - single job datum
type Job struct {
	ID            string `json:"id"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	ImageLocation string `json:"imageLocation"`
	Runner        string `json:"runner"`
}
