package service

import (
	"context"

	"weather-api/internal/model"
)

type TokenManager interface {
	Generate(user *model.User) (string, error)
}

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherData, error)
}

type GeoProvider interface {
	GetCoordinates(ctx context.Context, city, countryCode string) (float64, float64, error)
}

type CountryProvider interface {
	GetCountry(ctx context.Context, country string) (*CountryData, error)
	GetCities(ctx context.Context, country string) ([]string, error)
}

type WeatherLookupService interface {
	GetWeatherByCity(ctx context.Context, city string) (*model.WeatherResult, error)
}

type WeatherData struct {
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	Time        string
}

type CountryData struct {
	Name string
	Code string
}
