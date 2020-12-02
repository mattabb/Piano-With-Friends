package main

import (
	// Internal Dependencies
	"log"
	"net/http"
	"os"

	// External Dependencies
	"github.com/gorilla/mux"
)

func main() {

	log.Println("Starting the backend server at http://localhost:8000/")

	router := mux.NewRouter()

	// Whitelists our front-end
	router.Host("localhost:8080")

	// Add routes
	AddAppRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Port our back-end will connect to
	log.Fatal(http.ListenAndServe(":"+port, router))
}
