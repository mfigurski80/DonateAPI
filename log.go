package main

import (
	"fmt"
	"net/http"
	"time"
)

func log(s string) {
	fmt.Printf("[%s] %s\n", time.Now().Format("Jan 02 15:04:05"), s)
}

func logRequest(r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("Jan 02 15:04:05"), r.Method, r.URL)
}
