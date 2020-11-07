package handlers


// returns an instance of the pool (holds sockets)
func NewPool() *Pool {
	return &Pool{
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[*Client]bool),
	}
}


func (pool *Pool) Run() {
	for {
		select {
		case client := <- hub.register:
			HandleUserRegisterEvent(pool, client)

		case client := <- hub.unregister:
			HandleUserDisconnectEvent(pool, client)
		}
	}
}