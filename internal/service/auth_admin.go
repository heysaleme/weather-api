package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"weather-api/internal/errs"
	"weather-api/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) EnsureAdminAccount(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if email == "" && password == "" {
		return nil
	}
	if email == "" || password == "" {
		return errs.InvalidInput("bootstrap admin email and password are required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errs.InvalidInput("valid bootstrap admin email is required")
	}
	if len(password) < 12 {
		return errs.InvalidInput("bootstrap admin password must be at least 12 characters")
	}

	existing, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.Role != model.RoleAdmin {
			return errs.Conflict("bootstrap admin email already exists with non-admin role")
		}
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Email:        email,
		Role:         model.RoleAdmin,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	return s.users.Create(ctx, admin)
}
