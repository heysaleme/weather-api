package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

type LoggingMiddleware struct {
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := stringsOrFallback(strings.TrimSpace(r.Header.Get(requestIDHeader)), newRequestID())
		w.Header().Set(requestIDHeader, requestID)

		requestLogger := m.logger.With(
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		ctx := WithLogger(WithRequestID(r.Context(), requestID), requestLogger)
		r = r.WithContext(ctx)

		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		m.logger.Info(
			"http_request",
			zap.String("request_id", requestID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status_code", recorder.statusCode),
			zap.Duration("duration", duration),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func newRequestID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buf)
}

func stringsOrFallback(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
