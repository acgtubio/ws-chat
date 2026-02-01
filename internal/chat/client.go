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
	Id      string
	Conn    *websocket.Conn
	Receive chan Message // Channel for receiving data from the room.
	Logger  *zap.SugaredLogger
	rooms   map[string]*Room
}

// TODO: rooms should be based on the available rooms of the client. Add this functionality when introducing persistent mesesages.
func NewClient(id string, conn *websocket.Conn, logger *zap.SugaredLogger) *Client {
	return &Client{
		Id:      id,
		Conn:    conn,
		Logger:  logger,
		Receive: make(chan Message),
		rooms:   make(map[string]*Room),
	}
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

func (c *Client) unregister() {
	for _, room := range c.rooms {
		room.EmitEvent(RoomEvent{
			EventType: Leave,
			Member:    c,
		})
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.unregister()
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

		// TODO: Add additional author id validation later that is based on authentication.
		if msg.AuthorID == "" {
			c.Logger.Errorw("Empty author id.")
			break
		}

		// TODO: Add a custom message that will be received by the user. Not sure if connection should be dropped.
		if msg.Room == "" {
			c.Logger.Errorw("Empty room.")
			break
		}

		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now()
		}

		c.sendMessage(msg)
	}

	c.Logger.Debugw("closing read loop for client", "client", c.Id)
}

func (c *Client) JoinRoom(r *Room) {
	c.rooms[r.Id] = r
}

func (c *Client) sendMessage(msg Message) {
	room, ok := c.rooms[msg.Room]
	if !ok {
		return
	}

	room.Broadcast(msg)
}
