package api

import (
	"net/http"
	"time"

	"github.com/mfigurski80/DonateAPI/store"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixNano() / int64(time.Millisecond)
		next.ServeHTTP(w, r)
		end := time.Now().UnixNano() / int64(time.Millisecond)

		store.L.Printf("%s %s (+%vms)", r.Method, r.URL.Path, end-start)
	})
}
