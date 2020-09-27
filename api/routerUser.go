package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GET `/{userId}`
func getUser(w http.ResponseWriter, r *http.Request) {
	// auth
	userID := mux.Vars(r)["userId"]
	user, ok := authRequest(w, r)
	if !ok {
		return
	}
	if user.Username != userID {
		respondUnauthorized(w, "You are not authorized to view this page")
		return
	}

	// respond
	bytes, err := json.Marshal(user)
	if err != nil {
		respondInternalServerError(w, err)
		return
	}
	w.Write(bytes)
}

func addUserSubrouter(r *mux.Router) {
	userRouter := r.PathPrefix("/{userId}").Subrouter()
	userRouter.HandleFunc("", getUser).Methods(http.MethodGet)
}
