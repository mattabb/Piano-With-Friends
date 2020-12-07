package handlers

import (
	"bytes"
	"log"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

// Constants for delays, max message sizes and ping periods
const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 10000
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
		pool:                pool,
		webSocketConnection: connection,
		send:                make(chan SocketEventStruct),
		username:            username,
		recordNotes:         make(chan SocketEventStruct),
	}

	go client.writePump()
	go client.readPump()

	client.pool.register <- client
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
		EventName:    "join",
		EventPayload: EventPayloadStruct{User: client.username},
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
			EventName:    "disconnect",
			EventPayload: EventPayloadStruct{User: client.username},
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
func BroadcastSocketEventToAllClient(self *Client, payload SocketEventStruct) {

	for client := range self.pool.clients {
		if client != self {
			client.send <- payload
			log.Println("sent payload")
		}
		log.Print("the pool is: ", self.pool)
		log.Print("send channel ", client.send, "\n", "payload: ", payload)
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
	switch socketEventPayload.EventName {
	// When someone joins
	case "join":
		BroadcastSocketEventToAllClient(client, SocketEventStruct{
			EventName:    socketEventPayload.EventName,
			EventPayload: socketEventPayload.EventPayload,
		})

	// When someone disconnects
	case "disconnect":
		log.Print("Disconnect event triggered")
		BroadcastSocketEventToAllClient(client, SocketEventStruct{
			EventName:    socketEventPayload.EventName,
			EventPayload: socketEventPayload.EventPayload,
		})

	// When someone presses the keyboard
	case "keyboardPress":
		if client.recording {
			client.recordNotes <- socketEventPayload
		}
		BroadcastSocketEventToAllClient(client, socketEventPayload)

	// When someone presses record button
	case "record":
		beginRecord(client)
	}

}

/*
* @function readJSON
* @description
* Reads JSON from websocket, creates a JSON decoder and decodes data into payload

* @family Client
* @return socketEventPayload, error
 */
func (c *Client) readJSON() (SocketEventStruct, error) {
	var socketEventPayload SocketEventStruct

	_, payload, err := c.webSocketConnection.ReadMessage()
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoderErr := decoder.Decode(&socketEventPayload)

	log.Print("read JSON from ", c.username, " ... event name from websocket is: ", socketEventPayload.EventName)
	log.Print("read JSON from ", c.username, " ... event payload Username from websocket is: ", socketEventPayload.EventPayload.User)
	log.Print("read JSON from ", c.username, " ... event payload Message from websocket is: ", socketEventPayload.EventPayload.Message)
	log.Print("read JSON from ", c.username, " ... event payload Time from websocket is: ", socketEventPayload.EventPayload.Time)

	if decoderErr != nil {
		log.Printf("error: %v", decoderErr)
		return socketEventPayload, decoderErr
	}

	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("error ===: %v", err)
		}
		return socketEventPayload, err
	}

	return socketEventPayload, nil
}

/*
* @function encodeJSON
* @description
* Encodes JSON from send channel, and sets the write deadline

* @family Client
* @param payload SocketEventStruct
* @return []byte, error
 */
func (c *Client) encodeJSON(payload SocketEventStruct) ([]byte, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Print("error marshalling payload")
		return jsonPayload, err
	}

	// Write now
	c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))

	return jsonPayload, nil
}

/*
* @function writeJSON
* @description
* Creates JSON writer and writes to websocket using writerf

* @family Client
* @param payload []byte
* @return []byte, error
 */
func (c *Client) writeJSON(jsonData []byte) error {
	w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, errr := w.Write(jsonData)
	if errr != nil {
		log.Print("error when trying to write", errr)
		return errr
	}

	w.Close()
	// if err := w.Close(); err != nil{
	//	  return
	// }
	log.Println("Closed writer")

	return nil
}

/*
* @function readPump
* @description
* Reads data from websocket

* @family Client
* @return N/A
 */
func (c *Client) readPump() {
	// Read from websocket

	c.webSocketConnection.SetReadLimit(maxMessageSize)
	c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocketConnection.SetPongHandler(func(string) error { c.webSocketConnection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		socketEventPayload, err := c.readJSON()
		if err != nil {
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.webSocketConnection.Close()
	}()

	for {
		select {
		case payload, ok := <-c.send:

			log.Print("Hit writepump for ", c.username, " payload is: ", payload)
			jsonPayload, encodeErr := c.encodeJSON(payload)
			if encodeErr != nil {
				return
			}

			if !ok {
				log.Print("not ok")
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.writeJSON(jsonPayload)
			if err != nil {
				log.Println(err)
				return
			}

			return

		case <-ticker.C:
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
