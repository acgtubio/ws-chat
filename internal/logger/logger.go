package logger

import "go.uber.org/zap"

func NewLogger() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true

	logger, _ := cfg.Build()
	sugar := logger.Sugar()

	return sugar
}
