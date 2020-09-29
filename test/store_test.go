package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mfigurski80/DonateAPI/store"
)

func TestStore(t *testing.T) {
	store.Init("../data")

	users, err := store.ReadUsers()
	if err != nil {
		t.Errorf("store.ReadUsers() error =  %v", err)
	}

	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprint(rand.Uint64())
	user := store.User{
		Username: id,
	}
	users[id] = user

	err = store.WriteUsers(users)
	if err != nil {
		t.Errorf("store.WriteUsers() error = %v", err)
	}
	users, err = store.ReadUsers()
	if err != nil {
		t.Errorf("store.ReadUsers() error =  %v", err)
	}

	_, ok := users[id]
	if !ok {
		t.Errorf("store.ReadUsers() usermap does not contain user '%s' written by store.WriteUsers()", id)
	}

	delete(users, id)
	err = store.WriteUsers(users)
	if err != nil {
		t.Errorf("store.WriteUsers() error = %v", err)
	}

}
