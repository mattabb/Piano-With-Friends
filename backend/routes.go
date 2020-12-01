package main

import (
	// Internal Dependencies
	"log"
	"net/http"

	"github.com/mattabb/server/handlers"

	// External Dependencies
	"math/rand"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

/*
* @function setStaticFolder
* @description
* Sets static folder where our JS is held... not sure if we need this tbh

* @param {*mux.Router} route => mux router
* @return N/A
 */
func setStaticFolder(route *mux.Router) {
	// Serve all of our JS
	fs := http.FileServer(http.Dir("../UI/src"))
	route.PathPrefix("/UI/src").Handler(http.StripPrefix("/UI/src", fs))
}

/*
* @function addTrailingIntToUser
* @description
* Adds trailing integers to users so we don't have users with the same username

* @param string username
* @return string username1234
 */
func addTrailingIntToUser(username string) string {
	for i := 0; i < 4; i++ {
		randInt := rand.Intn(10)
		randIntString := strconv.Itoa(randInt)
		username += randIntString
	}
	return username
}

/*
* @function AddAppRoutes
* @description
* Adds app routes/endpoints to listen to for our server

* @exported: true
* @param {*mux.Router} route => mux router
* @return N/A
 */
func AddAppRoutes(route *mux.Router) {

	setStaticFolder(route)

	// Implement websockets and handlers
	pool := handlers.NewPool()
	go pool.Run()

	// Websocket handling
	route.HandleFunc("/ws/{username}", func(responseWriter http.ResponseWriter, request *http.Request) {

		var upgrader = websocket.Upgrader{
			// This allows connections from anywhere... the logic commented out below can restrict
			// origins, in our case we would want to restrict to our website homepage as we get further
			CheckOrigin: func(r *http.Request) bool { return true },
		}

		// if req.Header.Get("Origin") != "http://"+req.Host {
		// 	http.Error(w, "Origin not allowed", http.StatusForbidden)
		// 	return
		// }

		username := mux.Vars(request)["username"]

		username = addTrailingIntToUser(username)

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(pool, connection, username)
	}) // end of websocket handling

	log.Println("Routes loaded.")
}
