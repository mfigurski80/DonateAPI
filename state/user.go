package state

// User - single user datum
type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password uint32   `json:"password"`
	Authored []string `json:"authored"`
	Running  []string `json:"running"`
}
