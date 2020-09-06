package main

import "net/http"

func InternalServerError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err))
}

func UnsupportedMediaType(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Write([]byte("This `Content-Type` is not supported at this path"))
}

func BadRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("This resource could not be found"))
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("You are not authorized to access this resource"))
}
