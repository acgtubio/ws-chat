package main

import (
	"fmt"
	"net/http"

	"github.com/acgtubio/ws-chat/config"
	"github.com/acgtubio/ws-chat/internal/chat"
	"github.com/acgtubio/ws-chat/routes"
)

func main() {
	// ctx := context.Background()

	logger := NewLogger()
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

	routerDependencies := &routes.RouterDependencies{
		Logger: logger,
	}
	router, err := routes.SetupRoutes(routerDependencies)
	if err != nil {
		logger.Errorw("Error setting up routes.",
			"error", err,
		)
		return
	}

	// Starts a goroutine here for the hub.
	logger.Infow("Starting chat rooms.")
	chat.InitializeChat(logger)

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
