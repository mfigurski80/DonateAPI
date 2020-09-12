package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type userReader struct {
	path string
	sync.Mutex
	cache map[string]User
}

func newUsersReader(p string) *userReader {
	os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0666)

	return &userReader{
		path:  p,
		cache: map[string]User{},
	}
}

func (r *userReader) AuthUser(u string, p string) (User, bool) {
	users := r.Read()
	user, ok := users[u]
	if !ok {
		return User{}, false
	}
	if user.Password != HashPassword(p) {
		return User{}, false
	}
	return user, true
}

func (r *userReader) AuthRequest(req *http.Request) (User, bool) {
	username, pass, ok := req.BasicAuth()
	if !ok {
		return User{}, false
	}
	return r.AuthUser(username, pass)
}

func (r *userReader) Read() map[string]User {
	r.Lock()
	defer r.Unlock()
	if len(r.cache) > 0 {
		return r.cache
	}

	file, err := ioutil.ReadFile(r.path)
	if err != nil {
		panic(fmt.Sprintf("Error opening file '%s'", r.path))
	}

	var users map[string]User
	json.Unmarshal(file, &users)
	if len(users) == 0 {
		users = map[string]User{}
	}

	r.cache = users
	return users
}

func (r *userReader) Write(users map[string]User) {
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

// UserState - allows persisting User data
var UserState = newUsersReader("./data/Users.json")
