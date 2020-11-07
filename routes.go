package main


import (
	// Internal Dependencies
	"log"
	"net/http"
	// External Dependencies
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

)

func setStaticFolder(route *mux.Router) {
	fs := http.FileServer(http.Dir("./UI/public"))
	// route.PathPrefix("/")
}

func AddApproutes(route *mux.Router) {

	setStaticFolder(route)

	// Implement websockets and handlers
}