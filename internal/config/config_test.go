package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("BOOTSTRAP_ADMIN_EMAIL", "admin@example.com")
	t.Setenv("BOOTSTRAP_ADMIN_PASSWORD", "very-strong-pass")

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "9090", cfg.HTTPPort)
	assert.Equal(t, "secret", cfg.JWTSecret)
	assert.Equal(t, "admin@example.com", cfg.BootstrapAdminEmail)
	assert.Equal(t, "very-strong-pass", cfg.BootstrapAdminPass)
}

func TestLoadConfigRequiresJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	cfg, err := Load()
	require.Error(t, err)
	assert.Nil(t, cfg)
}
