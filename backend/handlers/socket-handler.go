package handlers


import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"encoding/json"
)


/*
* @function CreateNewSocketUser
* @description 
* Creates a new Socket User using the client struct
* calls writePump() and readPump() for the client

* @exported: true
* @param {*Pool} pool => Contains all of our clients
* @param {*websocket.Conn} => Connection to websocket
* @param String username => Clients username
* @return N/A
*/
func CreateNewSocketUser(pool *Pool, connection *websocket.Conn, username string) {
	client := &Client{
		pool: 					pool,
		webSocketConnection:	connection,
		send:					make(chan SocketEventStruct),
		username:				username,
	}

	client.pool.register <- client
	log.Println("Socket user created with username:", username)

	// Write to the websocket (This is gonna contain the logic where we write to the websocket)
	go client.writePump()
	// Read from the websocket (This is gonna contain the logic where we READ from the websocket)
	go client.readPump()
}

/*
* @function HandleUserRegisterEvent
* @description 
* Handler for when user registers
* calls handleSocketPayloadEvents

* @exported: true
* @param {*Pool} pool => Contains all of our clients
* @param {*Client} => Pointer to client
* @return N/A
*/
func HandleUserRegisterEvent(pool *Pool, client *Client) {
	pool.clients[client] = true
	handleSocketPayloadEvents(client, SocketEventStruct{
		EventName: "join",
		EventPayload: client.username,
	})
}

/*
* @function HandleUserDisconnectEvent
* @description 
* Handler for when user disconnects
* calls handleSocketPayloadEvents and deletes client from our pool

* @exported: true
* @param {*Pool} pool => Contains all of our clients
* @param {*Client} => Pointer to client
* @return N/A
*/
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


/*
* @function BroadcastSocketEventToAllClient
* @description 
* Broadcasts socket event to all available clients in pool by sending payload in client.send
* This will later be handled in readPump()

* @exported: true
* @param {*Pool} pool => Contains all of our clients
* @param {SocketEventStruct} payload => contains message being sent along websocket
* @return N/A
*/
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


/*
* @function handleSocketPayloadEvents
* @description 
* Handles the different events in our websocket (join, disconnect, keyboardPress)
* then sends them out to all clients using BroadcastSocketEventToAllClient

* @param {*Client} client => Our client
* @param {SocketEventStruct} socketEventPayload => contains message being sent along websocket
* @return N/A
*/
func handleSocketPayloadEvents(client *Client, socketEventPayload SocketEventStruct) {
	var socketEventResponse SocketEventStruct
	switch socketEventPayload.EventName {
	// When someone joins 
	case "join":
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
		socketEventResponse.EventName = "keyBdPressResponse"
		socketEventResponse.EventPayload = map[string]interface{}{
			"username": client.username,
			"message":	socketEventPayload.EventPayload,
		}
		BroadcastSocketEventToAllClient(client.pool, socketEventResponse)
	}
}

/*
* @function readPump
* @description 
* Reads data from websocket (currently configured for JSON... we want GOB if possible)

* @family Client
* @return N/A
*/
func (c *Client) readPump() {
	// Read from websocket
	var socketEventPayload SocketEventStruct

	defer unRegisterAndCloseConnection(c)

	for {
		_, payload, err := c.webSocketConnection.ReadMessage()
		
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoderErr := decoder.Decode(&socketEventPayload)
		log.Print(socketEventPayload)

		if decoderErr != nil {
			log.Printf("error: %v", decoderErr)
			break
		}

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error ===: %v", err)
			}
			break
		}
		handleSocketPayloadEvents(c, socketEventPayload)
	}
}

/*
* @function writePump
* @description 
* Writes data to websocket (currently configured for JSON... we want GOB if possible)

* @family Client
* @return N/A
*/
func (c *Client) writePump() {
	// Write to websocket
	// ticker := time.NewTicker(someDelay)
	// defer func() {
	// 	ticket.Stop()
	// 	c.webSocketConnection.Close()
	// }()
	for {
		select {
		case payload, ok := <- c.send:
			// Encode our payload
			reqBodyBytes := new(bytes.Buffer)
			json.NewEncoder(reqBodyBytes).Encode(payload)
			finalPayload := reqBodyBytes.Bytes()

			
			c.webSocketConnection.SetWriteDeadline(time.Now())
			if !ok {
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(finalPayload)

			n := len(c.send)
			for i := 0; i < n; i++ {
				json.NewEncoder(reqBodyBytes).Encode(<-c.send)
				w.Write(reqBodyBytes.Bytes())
			}

			if err := w.Close(); err != nil {
				return
			}

			
		}
	}
}

/*
* @function unRegisterAndCloseConnection
* @description 
* Unregisters client from pool and closes websocket

* @param {*Client} c => Client
* @return N/A
*/
func unRegisterAndCloseConnection(c *Client) {
	c.pool.unregister <- c
	c.webSocketConnection.Close()
}

/*
* @function setSocketPayloadReadConfig
* @description 
* Sets our configurations => message delay limits, message length limits, etc.

* @param {*Client} c => Contains client
* @return N/A
*/
func setSocketPayloadReadConfig(c *Client) {
	// Set all of our configurations... => Message delay limits, etc.
}