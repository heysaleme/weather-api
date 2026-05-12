package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type Config struct {
	HTTPPort            string
	JWTSecret           string
	JWTExpiration       time.Duration
	BootstrapAdminEmail string
	BootstrapAdminPass  string
}

func Load() (*Config, error) {
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if jwtSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &Config{
		HTTPPort:            defaultString(strings.TrimSpace(os.Getenv("HTTP_PORT")), "8080"),
		JWTSecret:           jwtSecret,
		JWTExpiration:       24 * time.Hour,
		BootstrapAdminEmail: strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_EMAIL")),
		BootstrapAdminPass:  strings.TrimSpace(os.Getenv("BOOTSTRAP_ADMIN_PASSWORD")),
	}, nil
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
