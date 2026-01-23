package routes

import (
	"github.com/acgtubio/ws-chat/handlers/chat"
	hub "github.com/acgtubio/ws-chat/internal/chat"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RouterDependencies struct {
	Logger *zap.SugaredLogger
	Hub    *hub.ChatHub
}

func SetupRoutes(dependencies *RouterDependencies) (*mux.Router, error) {
	mux := mux.NewRouter()

	mux.Path("/api/chat/{id}").
		Handler(
			chat.NewChatHandler(dependencies.Logger, dependencies.Hub),
		)

	return mux, nil
}
