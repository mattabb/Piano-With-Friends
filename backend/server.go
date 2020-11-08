package main

import (
	// Internal Dependencies
	"log"
	"net/http"

	// External Dependencies
	"github.com/gorilla/mux"
)

func main() {

	log.Println("Starting the backend server at http://localhost:8000/")

	router := mux.NewRouter()
	router.Host("localhost:8080")
	AddAppRoutes(router)

	log.Fatal(http.ListenAndServe(":8000", router))
}