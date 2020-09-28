package store

import (
	"io"
	"log"
	"net/http"
	"os"
)

var f, _ = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
var w = io.MultiWriter(os.Stdout, f)

// L is a logger that outputs to both stdout and logfile
var L = log.New(w, "", log.LstdFlags)

// Log - logs given string at time
func Log(s string) {
	L.Println(s)
}

// LogRequest - logs given request at time
func LogRequest(r *http.Request) {
	L.Printf("%s %s\n", r.Method, r.URL)
}
