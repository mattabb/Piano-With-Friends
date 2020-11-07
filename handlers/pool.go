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
		case client := <- pool.register:
			HandleUserRegisterEvent(pool, client)

		case client := <- pool.unregister:
			HandleUserDisconnectEvent(pool, client)
		}
	}
}