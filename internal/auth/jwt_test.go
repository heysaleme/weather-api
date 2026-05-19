package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weather-api/internal/model"
)

func TestJWTManagerGenerateAndParse(t *testing.T) {
	manager, err := NewJWTManager("test-secret", time.Hour)
	require.NoError(t, err)

	token, err := manager.Generate(&model.User{ID: 7, Email: "user@example.com", Role: model.RoleAdmin})
	require.NoError(t, err)

	claims, err := manager.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, int64(7), claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, model.RoleAdmin, claims.Role)
}

func TestJWTManagerRequiresSecret(t *testing.T) {
	manager, err := NewJWTManager("", time.Hour)
	require.Error(t, err)
	assert.Nil(t, manager)
}
