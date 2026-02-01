package chat

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

type RoomEventType string

const (
	Join  RoomEventType = "join"
	Leave RoomEventType = "leave"
)

type RoomEvent struct {
	EventType RoomEventType
	Member    *Client
}

type Room struct {
	Id        string
	members   map[*Client]bool
	broadcast chan Message
	roomEvent chan RoomEvent
	logger    *zap.SugaredLogger
}

func (r *Room) Run() {
	for {
		select {
		case message := <-r.broadcast:
			r.handleBroadcast(message)
		case event := <-r.roomEvent:
			r.handleEvent(event)
		}
	}
}

func (r *Room) handleEvent(event RoomEvent) {
	switch event.EventType {
	case Join:
		r.handleJoin(event)
	case Leave:
		r.handleLeave(event)
	}
}

func (r *Room) handleJoin(event RoomEvent) {
	// Client is added to the room, while room is also added to the client.
	r.members[event.Member] = true
	event.Member.JoinRoom(r)

	r.fanoutRoom(Message{
		Type:      MemberJoin,
		Timestamp: time.Now(),
		Payload:   fmt.Sprintf("%s has connected.", event.Member.Id),
		Room:      r.Id,
	})
}

func (r *Room) handleLeave(event RoomEvent) {
	memberId := event.Member.Id

	close(event.Member.Receive)
	delete(r.members, event.Member)

	// Send messages to members
	r.fanoutRoom(Message{
		Type:      MemberJoin,
		Timestamp: time.Now(),
		Payload:   fmt.Sprintf("%s has disconnected.", memberId),
	})

	if len(r.members) == 0 {
		return
	}
}

func (r *Room) handleBroadcast(msg Message) {
	switch msg.Type {
	case MemberChat:
		r.fanoutRoom(msg)
	default:
		// TODO: Add other custom room handling. For instance when a user joins a room while having a current session.
	}
}

func (r *Room) fanoutRoom(msg Message) {
	for member := range r.members {
		select {
		case member.Receive <- msg:
		default:
			close(member.Receive)
			delete(r.members, member)
		}
	}
}

func (r *Room) EmitEvent(event RoomEvent) {
	r.roomEvent <- event
}

func (r *Room) Broadcast(msg Message) {
	r.broadcast <- msg
}

func NewRoom(id string, logger *zap.SugaredLogger) *Room {
	// TODO: Add buffers
	return &Room{
		Id:        id,
		members:   make(map[*Client]bool),
		roomEvent: make(chan RoomEvent),
		broadcast: make(chan Message),
		logger:    logger,
	}
}
