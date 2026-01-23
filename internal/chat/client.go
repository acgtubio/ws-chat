package chat

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	writeWait  = 10 * time.Second
)

type Client struct {
	Id         string
	Conn       *websocket.Conn
	Receive    chan Message // Channel for receiving data from the hub.
	Send       chan<- Message
	Unregister chan<- *Client
	Logger     *zap.SugaredLogger
}

func (c *Client) WriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Receive:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			byteMsg, err := json.Marshal(message)
			if err != nil {
				return
			}
			w.Write(byteMsg)

			n := len(c.Receive)
			for range n {
				w.Write([]byte{'\n'})

				byteMsg, err := json.Marshal(<-c.Receive)
				if err != nil {
					return
				}
				w.Write(byteMsg)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.Logger.Errorw("Error reading message.",
				"error", err,
			)
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			c.Logger.Errorw("Error unmarshalling message.",
				"error", err,
			)
			break
		}

		// Send message to the broadcast in the hub.
		c.Send <- msg
	}
}
