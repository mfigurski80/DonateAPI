package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"sync"
)

type User struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password uint32   `json:"password"`
	Authored []string `json:"authored"`
	Running  []string `json:"running"`
}

func hashPassword(p string) uint32 {
	salt := os.Getenv("PASSWORD_SALT")
	if salt == "" {
		panic("env variable PASSWORD_SALT not set")
	}
	p = p + salt
	h := fnv.New32a()
	h.Write([]byte(p))
	return h.Sum32()
}

func authUser(u string, p string) (User, bool) {
	users := usersReader.read()
	user, ok := users[u]
	if !ok {
		return User{}, false
	}
	if user.Password != hashPassword(p) {
		return User{}, false
	}
	return user, true
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
		panic(fmt.Sprintf("Error marshaling json: %v", users))
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
