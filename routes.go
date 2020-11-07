package main


import (
	// Internal Dependencies
	"log"
	"net/http"
	"./handlers"
	// External Dependencies
	"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"

)

func setStaticFolder(route *mux.Router) {
	// Serve all of our JS
	fs := http.FileServer(http.Dir("./UI/src"))
	route.PathPrefix("/UI/src").Handler(http.StripPrefix("/UI/src", fs))
}

func AddAppRoutes(route *mux.Router) {

	setStaticFolder(route)

	// Implement websockets and handlers
	pool := handlers.NewPool()
	go pool.Run()
	log.Print("pool ran")
	// Set the default route to the index.html
	route.HandleFunc("/", handlers.RenderHome)
	
	// // Websocket handling
	// route.HandleFunc("/ws/{username}", func(responseWriter http.ResponseWriter, request *http.Request) {
		
	// 	var upgrader = websocket.Upgrader{
	// 		ReadBufferSize: 	1024,
	// 		WriteBufferSize: 	1024,
	// 	}

	// 	username := mux.Vars(request)["username"]

	// 	connection, err := upgrader.Upgrade(responseWriter, request, nil)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	handlers.CreateNewSocketUser(pool, connection, username)
	// })

	log.Println("Routes loaded.")
}
