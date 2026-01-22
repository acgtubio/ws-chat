package chat

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type chatHub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func (c *chatHub) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			delete(c.clients, client)
			close(client.send)
		// TODO: Remove broadcast
		case msg := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
}

func InitializeChat(
	logger *zap.SugaredLogger,
) {
	hub := &chatHub{
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	logger.Infow("Chat initializing.")
	go hub.run()
}
