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

type AuthService struct {
	users      UserStore
	jwtManager TokenManager
}

func NewAuthService(users UserStore, jwtManager TokenManager) *AuthService {
	return &AuthService{
		users:      users,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errs.InvalidInput("valid email is required")
	}
	if len(password) < 8 {
		return nil, errs.InvalidInput("password must be at least 8 characters")
	}

	existing, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errs.Conflict("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		Role:         model.RoleUser,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return sanitizeUser(user), nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return "", errs.InvalidInput("email and password are required")
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errs.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errs.Unauthorized("invalid credentials")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func sanitizeUser(user *model.User) *model.User {
	userCopy := *user
	userCopy.PasswordHash = ""
	return &userCopy
}
