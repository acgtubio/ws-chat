package main

import (
	"fmt"
	"net/http"

	"github.com/acgtubio/ws-chat/config"
	"github.com/acgtubio/ws-chat/internal/chat"
	"github.com/acgtubio/ws-chat/internal/logger"
	"github.com/acgtubio/ws-chat/routes"
)

func main() {
	// ctx := context.Background()

	logger := logger.NewLogger()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := logger.Sync()
		if err != nil {
			panic("Unable to flush logger on exit.")
		}
	}()

	hub := chat.NewChatHub()
	routerDependencies := &routes.RouterDependencies{
		Logger: logger,
		Hub:    hub,
	}
	router, err := routes.SetupRoutes(routerDependencies)
	if err != nil {
		logger.Errorw("Error setting up routes.",
			"error", err,
		)
		return
	}

	// Starts a goroutine here for the hub.
	chat.InitializeChat(logger, hub)

	logger.Infow("Chat service is running.",
		"port", cfg.Application.Port,
	)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Application.Port), router)

	if err != nil {
		logger.Errorw("Error starting application.",
			"error", err,
		)
		return
	}
}
