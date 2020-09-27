package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	AlteredID string `json:"alteredId"`
}

// base response makers

func makeResponse(success bool, message string, id string) []byte {
	resp, _ := json.Marshal(response{
		Success:   success,
		Message:   message,
		AlteredID: id,
	})
	return resp
}

func makeFailResponse(message string) []byte {
	return makeResponse(false, "Error: "+message, "")
}

func makeSuccessResponse(message string, id string) []byte {
	return makeResponse(true, "Success: "+message, id)
}

// functions that handle responseWriter:

func respondInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(makeFailResponse(fmt.Sprintf("internal error: %s", err.Error())))
}

func respondBadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(makeFailResponse(message))
}

func respondUnsupportedMediaType(w http.ResponseWriter, t string) {
	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Write(makeFailResponse(fmt.Sprintf("'%s' is not an accepted content type. Header 'Content-Type' should be 'application/json'", t)))
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(makeFailResponse(message))
}
