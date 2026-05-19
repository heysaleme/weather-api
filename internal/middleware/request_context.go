package middleware

import (
	"context"

	"go.uber.org/zap"
)

type requestContextKey string

const (
	loggerContextKey    requestContextKey = "request_logger"
	requestIDContextKey requestContextKey = "request_id"
)

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

func LoggerFromContext(ctx context.Context) (*zap.Logger, bool) {
	logger, ok := ctx.Value(loggerContextKey).(*zap.Logger)
	return logger, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	return requestID, ok
}
