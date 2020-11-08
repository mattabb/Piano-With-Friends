package handlers


import (
	"github.com/gorilla/websocket"
)

type SocketEventStruct struct {
	EventName string `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

// Client is a middleman between the websocket connection and the pool
type Client struct {
	pool 				*Pool
	webSocketConnection *websocket.Conn
	send 				chan SocketEventStruct
	username	 		string
}

// Maintains the set of active clients and broadcasts messages
type Pool struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}