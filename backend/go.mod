module github.com/mattabb/server

go 1.15

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
)

replace (
	github.com/mattabb/server/handlers => ./handlers
)
