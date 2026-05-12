package service

import (
	"context"
	"errors"
	"testing"

	"weather-api/internal/errs"
	"weather-api/internal/model"
	"weather-api/internal/repository"
)

type fakeTokenManager struct {
	token string
}

func (f fakeTokenManager) Generate(user *model.User) (string, error) {
	if f.token != "" {
		return f.token, nil
	}
	return "generated-token", nil
}

func TestAuthServiceRegisterCreatesRegularUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	service := NewAuthService(repo, fakeTokenManager{})

	user, err := service.Register(context.Background(), "USER@example.com", "strongpass123")
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.Role != model.RoleUser {
		t.Fatalf("expected role %q, got %q", model.RoleUser, user.Role)
	}
	if user.PasswordHash != "" {
		t.Fatalf("expected sanitized user without password hash")
	}

	stored, err := repo.GetByEmail(context.Background(), "user@example.com")
	if err != nil {
		t.Fatalf("GetByEmail returned error: %v", err)
	}
	if stored == nil || stored.PasswordHash == "" {
		t.Fatalf("expected stored user with password hash")
	}
}

func TestAuthServiceEnsureAdminAccountCreatesAdmin(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	service := NewAuthService(repo, fakeTokenManager{})

	if err := service.EnsureAdminAccount(context.Background(), "admin@example.com", "very-strong-pass"); err != nil {
		t.Fatalf("EnsureAdminAccount returned error: %v", err)
	}

	stored, err := repo.GetByEmail(context.Background(), "admin@example.com")
	if err != nil {
		t.Fatalf("GetByEmail returned error: %v", err)
	}
	if stored == nil {
		t.Fatalf("expected admin user to be created")
	}
	if stored.Role != model.RoleAdmin {
		t.Fatalf("expected role %q, got %q", model.RoleAdmin, stored.Role)
	}
}

func TestAuthServiceRegisterRejectsDuplicateEmail(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	service := NewAuthService(repo, fakeTokenManager{})

	if _, err := service.Register(context.Background(), "user@example.com", "strongpass123"); err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}

	_, err := service.Register(context.Background(), "user@example.com", "strongpass123")
	if err == nil {
		t.Fatal("expected duplicate registration error")
	}
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}
