package handlers


import (
	//"bytes"
	//"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Set a wait period before you send another keyboard action
	keyboardPressWait = 1 * time.Second
)


// CreateNewSocketUser creates a new socket user
func CreateNewSocketUser(pool *Pool, connection *websocket.Conn, username string) {
	client := &Client{
		pool: 					pool,
		webSocketConnection:	connection,
		send:					make(chan SocketEventStruct),
		username:				username,
	}

	client.pool.register <- client
	// Write to the websocket (This is gonna contain the logic where we write to the websocket)
	go client.writePump()
	// Read from the websocket (This is gonna contain the logic where we READ from the websocket)
	go client.readPump()
}

//HandleUserRegisterEvent will handle the Join event for new socket users
func HandleUserRegisterEvent(pool *Pool, client *Client) {
	pool.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName: "join",
		EventPayload: client.username,
	})
}

// HandleUserDisconnectEvent will handle the Disconnect event for socket users
func HandleUserDisconnectEvent(pool *Pool, client *Client) {
	_, ok := pool.clients[client]
	if ok {
		delete(pool.clients, client)
		close(client.send)

		handleSocketPayloadEvents(client, SocketEventStruct{
			EventName:		"disconnect",
			EventPayload:	client.username,
		}) 
	}
}

// BroadCastSocketEventToAllClient will emit the socket events to all socket users
func BroadcastSocketEventToAllClient(pool *Pool, payload SocketEventStruct) {
	for client := range pool.clients {
		select {
		case client.send <- payload:
		default:
			close(client.send)
			delete(pool.clients, client)
		}
	}
}

func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	var socketEventResponse SocketEventStruct
	switch socketEventPayload.EventName {
	// When someone joins 
	case "join":
		log.Printf("Join event triggered")
		BroadcastSocketEventToAllClient(client.pool, SocketEventStruct{
			EventName:		"join",
			EventPayload:	socketEventPayload.EventPayload,
		})
	// When someone disconnects
	case "disconnect:":
		log.Printf("Disconnect event triggered")
		BroadcastSocketEventToAllClient(client.pool, SocketEventStruct{
			EventName: 		"disconnect",
			EventPayload:	socketEventPayload.EventPayload,		
		})
	// When someone presses the keyboard
	case "keyboardPress":
		log.Printf("keyboard press event triggered")
		socketEventResponse.EventName = "keyboard press response"
		socketEventResponse.EventPayload = map[string]interface{}{
			"username": client.username,
			"message":	socketEventPayload.EventPayload,
		}
		BroadcastSocketEventToAllClient(client.pool, socketEventResponse)
	}
}

func (c *Client) readPump() {
	// Read from websocket
}


func (c *Client) writePump() {
	// Write to websocket
}

func unRegisterAndCloseConnection(c *Client) {
	c.pool.unregister <- c
	c.webSocketConnection.Close()
}

func setSocketPayloadReadConfig(c *Client) {
	// Set all of our configurations... => Message delay limits, etc.
}