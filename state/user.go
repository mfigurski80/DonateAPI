package state

// User - single user datum
type User struct {
	Username string   `json:"username"`
	Password uint32   `json:"password"`
	Authored []string `json:"authored"`
	Running  []string `json:"running"`
}
