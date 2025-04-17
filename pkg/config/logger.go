package config

import (
	"context"
	"log/slog"
)

type configCtxKey int

const configCtxLoggerKey configCtxKey = 1

// set logger as context value
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, configCtxLoggerKey, logger)
}

// extract logger from context
func LoggerFromContext(ctx context.Context) (*slog.Logger, bool) {
	logger, ok := ctx.Value(configCtxLoggerKey).(*slog.Logger)
	return logger, ok
}
