package chat

import (
	"net/http"

	"github.com/acgtubio/ws-chat/internal/chat"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type chatHandler struct {
	logger *zap.SugaredLogger
	hub    *chat.ChatHub
}

func NewChatHandler(
	logger *zap.SugaredLogger,
	hub *chat.ChatHub,
) http.Handler {
	return &chatHandler{
		logger,
		hub,
	}
}

func (d *chatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	// TODO: This is only a temporary code, used for testing.
	if userID == "" {
		d.logger.Warnw("Empty id, using temporary.")
		userID = uuid.New().String()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		d.logger.Errorw("Error upgrading to websocket.",
			"error", err,
		)
	}

	client := &chat.Client{
		Id:         userID,
		Conn:       conn,
		Send:       d.hub.BroadcastChan(),
		Receive:    make(chan chat.Message),
		Unregister: d.hub.UnregisterChan(),
		Logger:     d.logger,
	}

	d.hub.Register(client)

	go client.ReadLoop()
	go client.WriteLoop()
}
