package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	os.MkdirAll("./data", os.ModePerm)
	os.OpenFile("./data/Jobs.json", os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile("./data/Users.json", os.O_RDONLY|os.O_CREATE, 0666)

	r := mux.NewRouter()
	addAuthSubrouter(r)
	addUserSubrouter(r)
	addJobSubrouter(r)

	fmt.Println("Initialized. Serving on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
