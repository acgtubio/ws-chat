package chat

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type chatHandler struct {
	logger *zap.SugaredLogger
}

func NewChatHandler(logger *zap.SugaredLogger) http.Handler {
	return &chatHandler{
		logger,
	}
}

func (d *chatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		d.logger.Errorw("Error upgrading to websocket.",
			"error", err,
		)
	}
}
