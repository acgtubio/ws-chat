package chat

import "time"

type Message struct {
	Type      string
	Timestamp time.Time
	Payload   string
}
