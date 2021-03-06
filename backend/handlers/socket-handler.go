package handlers

import (
	"bytes"
	"log"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

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
		recordedNotes:       &[]SocketEventStruct{},
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
func BroadcastSocketEventToAllClient(self *Client, includeClient bool, payload SocketEventStruct) {

	for client := range self.pool.clients {
		if !(!includeClient && client == self) {
			client.send <- payload
		}
		// log.Print("the pool is: ", self.pool)
		// log.Print("send channel ", client.send, "\n", "payload: ", payload)
	}
	log.Print("we out")
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
		BroadcastSocketEventToAllClient(client, false, SocketEventStruct{
			EventName:    socketEventPayload.EventName,
			EventPayload: socketEventPayload.EventPayload,
		})

	// When someone disconnects
	case "disconnect":
		log.Print("Disconnect event triggered")
		BroadcastSocketEventToAllClient(client, false, SocketEventStruct{
			EventName:    socketEventPayload.EventName,
			EventPayload: socketEventPayload.EventPayload,
		})

	// When someone presses the keyboard
	case "keyboardPress":
		if client.recording {
			client.recordNotes <- socketEventPayload
		}
		BroadcastSocketEventToAllClient(client, false, SocketEventStruct{
			EventName:    socketEventPayload.EventName,
			EventPayload: socketEventPayload.EventPayload,
		})

	// When someone presses record button
	case "recordStart":
		if !client.recording {
			go beginRecord(client)
		}

	// When someone presses the record button while recording
	case "recordStop":
		client.recording = false

	case "recordPlay":
		go playRecording(client)
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
* @function readPump
* @description
* Reads data from websocket

* @family Client
* @return N/A
 */
func (c *Client) readPump() {

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
		unRegisterAndCloseConnection(c)
	}()

	for {
		select {
		case payload, ok := <-c.send:

			log.Print("Hit writepump for ", c.username, " payload is: ", payload)
			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				log.Print("error marshalling payload")
				return
			}

			// Write now
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				log.Print("not ok")
				c.webSocketConnection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.webSocketConnection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, errr := w.Write(jsonPayload)
			if errr != nil {
				log.Print("error when trying to write", errr)
				return
			}

			if err := w.Close(); err != nil {
				log.Print("closing the writer")
				return
			}

		case <-ticker.C:
			c.webSocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.webSocketConnection.WriteMessage(websocket.PingMessage, nil); err != nil {
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
	// close(c.send)
	c.webSocketConnection.Close()
}
