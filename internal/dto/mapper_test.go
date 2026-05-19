package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weather-api/internal/model"
)

func TestMappers(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)

	user := ToUserResponse(&model.User{ID: 1, Email: "user@example.com", Role: model.RoleUser, CreatedAt: now})
	require.NotNil(t, user)
	assert.Equal(t, int64(1), user.ID)

	city := ToCityResponse(&model.City{ID: 2, Name: "Almaty", CreatedAt: now})
	require.NotNil(t, city)
	assert.Equal(t, "Almaty", city.Name)

	weather := ToWeatherResponse(&model.WeatherResult{City: "Almaty", Temperature: 21})
	require.NotNil(t, weather)
	assert.Equal(t, "Almaty", weather.City)

	auth := ToAuthResponse("token")
	require.NotNil(t, auth)
	assert.Equal(t, "token", auth.AccessToken)
}
