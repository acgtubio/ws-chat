package routes

import (
	"github.com/acgtubio/ws-chat/handlers/chat"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type RouterDependencies struct {
	Logger *zap.SugaredLogger
}

func SetupRoutes(dependencies *RouterDependencies) (*mux.Router, error) {
	mux := mux.NewRouter()

	mux.Path("/api/chat").
		Handler(
			chat.NewChatHandler(dependencies.Logger),
		)

	return nil, nil
}
