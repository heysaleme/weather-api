package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddlewareAddsRequestIDAndLogsRequest(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	mw := NewLoggingMiddleware(logger)

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	req.Header.Set("X-Request-ID", "req-123")
	rr := httptest.NewRecorder()

	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := RequestIDFromContext(r.Context())
		require.True(t, ok)
		assert.Equal(t, "req-123", requestID)

		requestLogger, ok := LoggerFromContext(r.Context())
		require.True(t, ok)
		require.NotNil(t, requestLogger)

		w.WriteHeader(http.StatusAccepted)
	})).ServeHTTP(rr, req)

	require.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, "req-123", rr.Header().Get("X-Request-ID"))

	entries := recorded.All()
	require.Len(t, entries, 1)
	assert.Equal(t, "http_request", entries[0].Message)
	assert.Equal(t, "req-123", entries[0].ContextMap()["request_id"])
	assert.Equal(t, "/cities", entries[0].ContextMap()["path"])
	assert.Equal(t, int64(http.StatusAccepted), entries[0].ContextMap()["status_code"])
}
