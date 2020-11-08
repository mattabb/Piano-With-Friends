package handlers

import (
	"log"
)

/*
* @function NewPool
* @description 
* Creates a new pool of clients
* register: channel to send which client is registering
* unregister: channel to send which client is unregistering
* clients: channel to a map of clients with boolean values => tells us whether they are connected

* @exported: true
* @return {Pool} 
*/
func NewPool() *Pool {
	return &Pool{
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}

/*
* @function Run
* @description 
* Handles register and unregister events to our pool

* @family: Pool
* @exported: true
* @return N/A
*/
func (pool *Pool) Run() {
	for {
		select {
		case client := <- pool.register:
			HandleUserRegisterEvent(pool, client)
			log.Print(
				"event: ", "register \n",
				"pool: ", &pool.clients, "\n",
				"client: ", client, "\n",
			)
		case client := <- pool.unregister:
			HandleUserDisconnectEvent(pool, client)
			log.Print(
				"event: ", "unregister \n",
				"pool: ", &pool.clients, "\n",
				"client: ", client, "\n",
			)
		}
	}
}