package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"weather-api/internal/auth"
	"weather-api/internal/model"
)

type UserReader interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type contextKey string

const userContextKey contextKey = "auth_user"

type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	users      UserReader
}

func NewAuthMiddleware(jwtManager *auth.JWTManager, users UserReader) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		users:      users,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if header == "" {
			writeJSONError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			writeJSONError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		claims, err := m.jwtManager.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		storedUser, err := m.users.GetByID(r.Context(), claims.UserID)
		if err != nil || storedUser == nil {
			writeJSONError(w, http.StatusUnauthorized, "user no longer exists")
			return
		}

		user := &model.AuthUser{
			ID:    storedUser.ID,
			Email: storedUser.Email,
			Role:  storedUser.Role,
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok || user.Role != role {
				writeJSONError(w, http.StatusForbidden, "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) (*model.AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(*model.AuthUser)
	return user, ok
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
