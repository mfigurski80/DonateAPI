package types

// User is a single user datum
type User struct {
	Username string   `json:"username"`
	Password uint32   `json:"password"`
	Authored []string `json:"authored"`
	Running  []string `json:"running"`
}

// PostUserStruct is data type recieved when creating new user
type PostUserStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
