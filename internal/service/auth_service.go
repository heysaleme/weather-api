package service

import (
	"context"
	"net/mail"
	"strings"
	"time"

	"weather-api/internal/auth"
	"weather-api/internal/errs"
	"weather-api/internal/model"
	"weather-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users       repository.UserRepository
	jwtManager  *auth.JWTManager
	adminEmails map[string]struct{}
}

func NewAuthService(users repository.UserRepository, jwtManager *auth.JWTManager, adminEmails []string) *AuthService {
	adminSet := make(map[string]struct{}, len(adminEmails))
	for _, email := range adminEmails {
		key := strings.ToLower(strings.TrimSpace(email))
		if key != "" {
			adminSet[key] = struct{}{}
		}
	}

	return &AuthService{
		users:       users,
		jwtManager:  jwtManager,
		adminEmails: adminSet,
	}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)

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

	role := model.RoleUser
	if _, ok := s.adminEmails[email]; ok {
		role = model.RoleAdmin
	}

	user := &model.User{
		Email:        email,
		Role:         role,
		PasswordHash: string(hash),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return sanitizeUser(user), nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := strings.TrimSpace(req.Password)

	if email == "" || password == "" {
		return nil, errs.InvalidInput("email and password are required")
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.Unauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errs.Unauthorized("invalid credentials")
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{AccessToken: token}, nil
}

func sanitizeUser(user *model.User) *model.User {
	userCopy := *user
	userCopy.PasswordHash = ""
	return &userCopy
}
