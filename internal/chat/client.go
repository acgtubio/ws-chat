package chat

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Id         string
	Conn       *websocket.Conn
	Receive    <-chan Message
	Send       chan<- Message
	unregister chan *Client
	logger     *zap.SugaredLogger
}

func (c *Client) WriteLoop() {
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.logger.Errorw("Error reading message.",
				"error", err,
			)
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			c.logger.Errorw("Error unmarshalling message.",
				"error", err,
			)
			break
		}

		// Send message to the broadcast in the hub.
		c.Send <- msg
	}
}
