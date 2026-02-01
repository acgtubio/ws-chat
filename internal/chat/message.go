package chat

import "time"

type MessageType string

const (
	MemberLeave = "leave"
	MemberJoin  = "join"
	MemberChat  = "msg"
)

type Message struct {
	Type      MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   string      `json:"payload"`
	AuthorID  string      `json:"authorId"`
	Room      string      `json:"room"`
}
