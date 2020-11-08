package main


import (
	// Internal Dependencies
	"log"
	"net/http"
	"./handlers"
	// External Dependencies
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"math/rand"
	"strconv"

)

func setStaticFolder(route *mux.Router) {
	// Serve all of our JS
	fs := http.FileServer(http.Dir("../UI/src"))
	route.PathPrefix("/UI/src").Handler(http.StripPrefix("/UI/src", fs))
}

func AddAppRoutes(route *mux.Router) {

	setStaticFolder(route)

	// Implement websockets and handlers
	pool := handlers.NewPool()
	go pool.Run()
	log.Print("pool ran")
	
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

		// Generate some random numbers to append to username so we don't have overlapping usernames
		for i := 0; i < 4; i++ {
			randInt := rand.Intn(10)
			randIntString := strconv.Itoa(randInt) 
			username += randIntString
		}

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(pool, connection, username)
	})

	log.Println("Routes loaded.")
}
