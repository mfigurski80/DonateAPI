package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/mfigurski80/DonateAPI/api"
)

var apiTests = []struct {
	method            string
	path              string
	username          string
	password          string
	contentType       string
	data              string
	expectStatus      int
	expectContentType string
	expectContent     string
}{
	{"POST", "/register", "", "", "application/json", `{"username":"TESTUSER", "password":"P"}`,
		http.StatusOK, "application/json", `^{.*"success": ?true.*alteredId": ?"TESTUSER".*}`},
}

var client = http.Client{
	Timeout: time.Duration(2 * time.Second),
}

func TestAPI(t *testing.T) {
	go api.Start(":8080")
	time.Sleep(2 * time.Second) // let api start up
	baseURL := "http://localhost:8080"

	for i, test := range apiTests {
		// set up request
		req, err := http.NewRequest(test.method, baseURL+test.path, bytes.NewBuffer([]byte(test.data)))
		if err != nil {
			t.Fatalf("Error making request for test #%d: %v", i, err)
		}

		req.Header.Set("Content-Type", test.contentType)
		req.SetBasicAuth(test.username, test.password)

		// do request
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error sending request for test#%d: %v", i, err)
		}

		// check answers
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error reading respond body for test#%d: %v", i, err)
		}
		body := string(bodyBytes)

		if resp.StatusCode != test.expectStatus {
			t.Fatalf("Test #%d got unexpected status %d (expected %d): \n\t%s", i, resp.StatusCode, test.expectStatus, body)
		}
		if ct := resp.Header.Get("Content-Type"); ct != test.expectContentType {
			t.Fatalf("Test #%d got unexpected content type %s (expected %s): \n\t%s", i, ct, test.expectContentType, body)
		}
		matched, err := regexp.MatchString(test.expectContent, body)
		if err != nil {
			panic(err)
		}
		if !matched {
			t.Fatalf("Test #%d could not match regex '%s' to response: \n\t'%s'", i, test.expectContent, body)
		}

	}
}
