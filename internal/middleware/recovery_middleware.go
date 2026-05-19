package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type RecoveryMiddleware struct {
	logger *zap.Logger
}

func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{logger: logger}
}

func (m *RecoveryMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger := m.logger
				if requestLogger, ok := LoggerFromContext(r.Context()); ok {
					logger = requestLogger
				}

				logger.Error("panic recovered", zap.String("panic", fmt.Sprint(rec)))
				writeJSONError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
