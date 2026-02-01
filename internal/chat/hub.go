package chat

import (
	"go.uber.org/zap"
)

type ChatHub struct {
	clients map[*Client]bool
	rooms   map[string]*Room
	event   chan HubEvent
	logger  *zap.SugaredLogger
}

type HubEvent struct {
	Client *Client
	Room   string
}

// TODO: Add handling for removing of inactive rooms.
func (c *ChatHub) run() {
	for {
		select {
		case event := <-c.event:
			c.subscribeToRoom(event)
		}
	}
}

func (c *ChatHub) subscribeToRoom(event HubEvent) {
	room, ok := c.rooms[event.Room]
	if !ok {
		room = NewRoom(event.Room, c.logger)
		c.rooms[event.Room] = room
		go room.Run()
	}

	joinEvent := RoomEvent{
		EventType: Join,
		Member:    event.Client,
	}

	room.EmitEvent(joinEvent)
}

func (c *ChatHub) EmitHubEvent(event HubEvent) {
	c.event <- event
}

func NewChatHub(logger *zap.SugaredLogger) *ChatHub {
	return &ChatHub{
		clients: make(map[*Client]bool),
		rooms:   make(map[string]*Room),
		event:   make(chan HubEvent),
		logger:  logger,
	}
}

func InitializeChat(
	logger *zap.SugaredLogger,
	hub *ChatHub,
) {
	logger.Infow("Chat initializing.")
	go hub.run()
}
