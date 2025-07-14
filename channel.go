package main

func newChannel(name string, channelType ChannelType) *Channel {
	return &Channel{
		name:        name,
		channelType: channelType,
		clients:     make(map[*Client]bool),
		broadcast:   make(chan []byte),
		shutdown:    make(chan bool),
	}
}

func (c *Channel) run(hubShutdown chan bool) {
	for {
		select {
		case <-c.shutdown:
			return
		case <-hubShutdown:
			return
		case message := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}