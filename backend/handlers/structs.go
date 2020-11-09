package handlers


import (
	"github.com/gorilla/websocket"
)

/*
	EventName: ('join'/'disconnect'/'keyBdPressResponse') JSON which contains the eventname
	EventPayload: JSON interface which contains our payload
*/
type SocketEventStruct struct {
	EventName string `json:"eventName"`
	EventPayload interface{} `json:"eventPayload"`
}

/* 
	Client is a middleman between the websocket connection and the pool
	pool: *Pool that client is associated with
	webSocketConnection: connection to websocket
	send: channel we are sending socket event to
	username: username associated to client
*/
type Client struct {
	pool 				*Pool
	webSocketConnection *websocket.Conn
	send 				chan SocketEventStruct
	username	 		string
}

/* 
	Maintains the set of active clients and broadcasts messages
	clients: map of clients with corresponding boolean values to tell whether connected
	broadcoat: channel to broadcoast bytes
	register: channel for clients that are registering
	unregister: channel for clients that are unregistering
*/
type Pool struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}