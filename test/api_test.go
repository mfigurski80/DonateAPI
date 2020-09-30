package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/mfigurski80/DonateAPI/api"
)

var apiTests = []struct {
	method        string
	path          string
	username      string
	password      string
	data          string
	expectStatus  int
	expectContent string
}{
	{"POST", "/register", "", "", `{"username":"TESTUSER", "password":"A"}`,
		http.StatusOK, `^{.*"success": ?true.*"alteredId": ?"TESTUSER".*}`},
	{"GET", "/TESTUSER", "TESTUSER", "A", "",
		http.StatusOK, `^{.*"username": ?"TESTUSER".*"authored": ?{}.*"running":.*}`},
	{"GET", "/job", "", "", "",
		http.StatusOK, `^\[\]`},
	{"POST", "/job", "TESTUSER", "A", `{"title":"TESTJOB", "description":"TESTDESC", "image": "TESTIMG"}`,
		http.StatusOK, `^{.*"success": ?true.*"alteredId": ?"TESTJOB".*}`},
	{"GET", "/job", "", "", "",
		http.StatusOK, `^\[.*{.*"title": ?"TESTJOB".*"author": ?"TESTUSER".*}.*\]`},
	{"GET", "/TESTUSER", "TESTUSER", "A", "",
		http.StatusOK, `^{.*"username": ?"TESTUSER".*"authored": ?{.*"TESTJOB":.*}.*}`},
	{"GET", "/TESTUSER/TESTJOB", "", "", "",
		http.StatusOK, `^{.*"title": ?"TESTJOB".*}`},
	{"DELETE", "/TESTUSER/TESTJOB", "TESTUSER", "A", "",
		http.StatusOK, `^{.*"success": ?true.*}`},
	{"DELETE", "/TESTUSER", "TESTUSER", "A", "",
		http.StatusOK, `^{.*"success": ?true.*}`},
}

var client = http.Client{
	Timeout: time.Duration(1 * time.Second),
}

func TestAPI(t *testing.T) {
	// clean from previous run
	err := os.RemoveAll("./data")
	if err != nil {
		panic(err)
	}
	// start api
	os.Setenv("PASSWORD_SALT", "salt")
	go api.Start(":8080")
	time.Sleep(2 * time.Second) // let api start up
	baseURL := "http://localhost:8080"

	for i, test := range apiTests {
		// set up request
		req, err := http.NewRequest(test.method, baseURL+test.path, bytes.NewBuffer([]byte(test.data)))
		if err != nil {
			t.Fatalf("Error making request for test %d: %v", i, err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(test.username, test.password)

		// do request
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error sending request for test %d: %v", i, err)
		}

		// check answers
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading respond body for test %d: %v", i, err)
		}
		body := string(bodyBytes)

		if resp.StatusCode != test.expectStatus {
			t.Fatalf("Test %d got unexpected status %d (expected %d): \n\t%s", i, resp.StatusCode, test.expectStatus, body)
		}
		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Test %d got unexpected content type %s (expected application/json): \n\t%s", i, ct, body)
		}
		matched, err := regexp.MatchString(test.expectContent, body)
		if err != nil {
			panic(err)
		}
		if !matched {
			t.Fatalf("Test %d could not match regex '%s' to response: \n\t'%s'", i, test.expectContent, body)
		}

	}
}
