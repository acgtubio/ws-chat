package chat

import (
	"go.uber.org/zap"
)

type ChatHub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func (c *ChatHub) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			delete(c.clients, client)
			close(client.Send)
		// TODO: Remove broadcast
		case msg := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(c.clients, client)
				}
			}
		}
	}
}

// BroadcastChan returns a channel for receiving messages.
// This is used to broadcast to all clients.
func (c *ChatHub) BroadcastChan() <-chan Message {
	return c.broadcast
}

func (c *ChatHub) Register(client *Client) {
	c.register <- client
}

func (c *ChatHub) Unregister(client *Client) {
	c.unregister <- client
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		broadcast:  make(chan Message),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func InitializeChat(
	logger *zap.SugaredLogger,
	hub *ChatHub,
) {
	logger.Infow("Chat initializing.")
	go hub.run()
}
