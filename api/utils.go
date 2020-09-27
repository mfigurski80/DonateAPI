package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mfigurski80/DonateAPI/store"
)

func unmarshalRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		respondBadRequest(w, "could not read request body")
		return false
	}
	if r.Header.Get("Content-Type") != "application/json" {
		respondUnsupportedMediaType(w, r.Header.Get("Content-Type"))
		return false
	}

	err = json.Unmarshal(bodyBytes, v)
	if err != nil {
		respondBadRequest(w, "request body is not valid json: "+err.Error())
		return false
	}

	return true
}

func authRequest(w http.ResponseWriter, r *http.Request) (store.User, bool) {
	username, pass, ok := r.BasicAuth()
	if !ok {
		respondUnauthorized(w, "Request does not have any authentication")
		return store.User{}, false
	}
	users, err := store.ReadUsers()
	if err != nil {
		respondBadRequest(w, err.Error())
		return store.User{}, false
	}

	user, ok := (*users)[username]
	if !ok {
		respondUnauthorized(w, fmt.Sprintf("User '%s' does not exist", username))
		return store.User{}, false
	}
	if user.Password != store.HashPassword(pass) {
		respondUnauthorized(w, fmt.Sprintf("Password '%s' does not match", pass))
		return store.User{}, false
	}

	return user, true
}
