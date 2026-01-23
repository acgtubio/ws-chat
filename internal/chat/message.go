package chat

import "time"

type Message struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   string    `json:"payload"`
}

// type Message string
