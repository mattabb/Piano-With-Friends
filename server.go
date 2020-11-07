package main

import (
	// Internal Dependencies
	"log"
	"net/http"
	// External Dependencies
	"github.com/gorilla/mux"
)

func main () {

	log.Println("Starting the backend server at http://localhost:8080/")

	route := mux.NewRouter()

	AddAppRoutes(route)

	log.Fatal(http.ListenAndServe(":8000", route))
}