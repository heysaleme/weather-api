package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather-api/internal/auth"
	"weather-api/internal/model"
)

type stubUserReader struct {
	user *model.User
	err  error
}

func (s stubUserReader) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.user, s.err
}

func TestAuthMiddlewareUsesStoredUserRole(t *testing.T) {
	jwtManager, err := auth.NewJWTManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}

	token, err := jwtManager.Generate(&model.User{
		ID:    7,
		Email: "user@example.com",
		Role:  model.RoleUser,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	mw := NewAuthMiddleware(jwtManager, stubUserReader{
		user: &model.User{
			ID:    7,
			Email: "user@example.com",
			Role:  model.RoleAdmin,
		},
	})

	next := RequireRole(model.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			t.Fatal("expected user in context")
		}
		if user.Role != model.RoleAdmin {
			t.Fatalf("expected stored admin role, got %q", user.Role)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	mw.Handle(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestAuthMiddlewareRejectsDeletedUser(t *testing.T) {
	jwtManager, err := auth.NewJWTManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewJWTManager returned error: %v", err)
	}

	token, err := jwtManager.Generate(&model.User{
		ID:    9,
		Email: "user@example.com",
		Role:  model.RoleUser,
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	mw := NewAuthMiddleware(jwtManager, stubUserReader{})
	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
