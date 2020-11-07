package main


import (
	// Internal Dependencies
	"log"
	"net/http"
	"./handlers"
	// External Dependencies
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

)


// idk how to serve VueJS here
func setStaticFolder(route *mux.Router) {
	fs := http.FileServer(http.Dir("./UI/public"))
	route.PathPrefix("/UI/public").Handler(http.StripPrefix("/UI/public/", fs))
}

func AddAppRoutes(route *mux.Router) {

	setStaticFolder(route)

	// Implement websockets and handlers
	pool := handlers.NewPool()
	go pool.Run()

	route.HandleFunc("/", handlers.RenderHome)

	route.HandleFunc("/ws/{username}", func(responseWriter http.ResponseWriter, request *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize: 	1024,
			WriteBufferSize: 	1024,
		}

		username := mux.Vars(request)["username"]

		connection, err := upgrader.Upgrade(responseWriter, request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		handlers.CreateNewSocketUser(pool, connection, username)
	})

	log.Println("Routes loaded.")
}
