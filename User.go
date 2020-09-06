package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type User struct {
	username string   `json:"username"`
	email    string   `json:"email"`
	password string   `json:"password"`
	authored []string `json:"authored"`
	running  []string `json:"running"`
}

type UserReader struct {
	path  string
	cache map[string]User
	sync.Mutex
}

func newUsersReader(p string) *UserReader {
	return &UserReader{
		path:  p,
		cache: map[string]User{},
	}
}

func (r *UserReader) read() map[string]User {
	r.Lock()
	if len(r.cache) > 0 {
		return r.cache
	}

	file, err := ioutil.ReadFile(r.path)
	if err != nil {
		panic(fmt.Sprintf("Error opening file '%s'", r.path))
	}

	var users map[string]User
	json.Unmarshal(file, &users)

	r.cache = users
	r.Unlock()
	return users
}

func (r *UserReader) write(users map[string]User) {
	file, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		panic(fmt.Sprintf("Error marshaling json: %s", users))
	}

	err = ioutil.WriteFile(r.path, file, 0644)
	if err != nil {
		panic(fmt.Sprintf("Error writing to file '%s'", r.path))
	}

	r.Lock()
	r.cache = users
	r.Unlock()
}

var usersReader = newUsersReader("./data/Users.json")
