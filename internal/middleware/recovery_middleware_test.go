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

func TestRecoveryMiddlewareRecoversAndLogsPanic(t *testing.T) {
	core, recorded := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	req = req.WithContext(WithLogger(WithRequestID(req.Context(), "req-999"), logger))
	rr := httptest.NewRecorder()

	NewRecoveryMiddleware(logger).Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})).ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"error":"internal server error"}`, rr.Body.String())

	entries := recorded.All()
	require.Len(t, entries, 1)
	assert.Equal(t, "panic recovered", entries[0].Message)
	assert.Equal(t, "boom", entries[0].ContextMap()["panic"])
}
