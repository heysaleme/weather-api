package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"weather-api/internal/errs"
)

func TestGeoClientRejectsEmptyCity(t *testing.T) {
	client := NewGeoClient(&http.Client{})

	lat, lon, err := client.GetCoordinates(context.Background(), "   ", "")
	require.Error(t, err)
	assert.Equal(t, 0.0, lat)
	assert.Equal(t, 0.0, lon)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)
}

func TestCountryClientRejectsEmptyCountry(t *testing.T) {
	client := NewCountryClient(&http.Client{})

	info, err := client.GetCountry(context.Background(), "   ")
	require.Error(t, err)
	assert.Nil(t, info)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)
}

func TestWeatherClientRejectsInvalidBaseURL(t *testing.T) {
	client := NewWeatherClient(&http.Client{})
	client.baseURL = ":"

	result, err := client.GetCurrentWeather(context.Background(), 1, 2)
	require.Error(t, err)
	assert.Nil(t, result)
}
