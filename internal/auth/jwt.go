package auth

import (
	"errors"
	"fmt"
	"time"

	"weather-api/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, expiration time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	return &JWTManager{
		secret:     []byte(secret),
		expiration: expiration,
	}, nil
}

func (m *JWTManager) Generate(user *model.User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Method.Alg())
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
