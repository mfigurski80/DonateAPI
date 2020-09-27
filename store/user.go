package store

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

// JobReference is absolute reference to specific job
type JobReference struct {
	User  string `json:"user"`
	Title string `json:"title"`
}

// User is a single user datum
type User struct {
	Username string         `json:"username"`
	Password uint32         `json:"password"`
	Authored JobMap         `json:"authored"`
	Running  []JobReference `json:"running"`
}

// UserMap is a map of usernames to user data
type UserMap map[string]User

var userPath = "./data/users.json"

var userCache = struct {
	c UserMap
	sync.Mutex
}{
	c: UserMap{},
}

// ReadUsers reads all users from file
func ReadUsers() (UserMap, error) {
	userCache.Lock()
	defer userCache.Unlock()
	if len(userCache.c) > 0 {
		return userCache.c, nil
	}

	file, err := ioutil.ReadFile(userPath)
	if err != nil {
		panic(err)
	}

	var users UserMap
	err = json.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}

	userCache.c = users
	return users, nil
}

// WriteUsers writes given user map to file
func WriteUsers(users UserMap) error {
	file, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(userPath, file, 0644)
	if err != nil {
		panic(err)
	}

	userCache.Lock()
	userCache.c = users
	userCache.Unlock()

	return nil
}
