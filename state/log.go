package state

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func appendToLog(s string) {
	f, err := os.OpenFile("./data/main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(s)
}

func printTime() {
	fmt.Printf("[%s] ", time.Now().Format("Jan 02 15:04:05"))
}

// Log - logs given string at time
func Log(s string) {
	appendToLog(s)
	printTime()
	fmt.Println(s)
}

// LogRequest - logs given request at time
func LogRequest(r *http.Request) {
	log := fmt.Sprintf("%s %s", r.Method, r.URL)
	appendToLog(log)
	printTime()
	fmt.Println(log)
}
