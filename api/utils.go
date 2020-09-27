package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
