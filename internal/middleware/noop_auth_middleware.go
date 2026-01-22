package middleware

import (
	"net/http"

	"github.com/acgtubio/ws-chat/config"
	"go.uber.org/zap"
)

type NoOpAuth struct {
	cfg    config.Config
	logger *zap.SugaredLogger
}

func NewNoOpAuth(cfg config.Config, logger *zap.SugaredLogger) *NoOpAuth {
	return &NoOpAuth{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *NoOpAuth) NewAuthMiddlware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
