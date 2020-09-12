package state

import (
	"fmt"
	"net/http"
	"time"
)

// Log - logs given string at time
func Log(s string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("Jan 02 15:04:05"), s)
}

// LogRequest - logs given request at time
func LogRequest(r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("Jan 02 15:04:05"), r.Method, r.URL)
}
