package client

import (
	"context"

	"weather-api/internal/service"
)

type WeatherProviderAdapter struct {
	Client *WeatherClient
}

func (a WeatherProviderAdapter) GetCurrentWeather(ctx context.Context, lat, lon float64) (*service.WeatherData, error) {
	weather, err := a.Client.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, err
	}

	return &service.WeatherData{
		Temperature: weather.Temperature,
		WindSpeed:   weather.WindSpeed,
		WeatherCode: weather.WeatherCode,
		Time:        weather.Time,
	}, nil
}

type CountryProviderAdapter struct {
	Client *CountryClient
}

func (a CountryProviderAdapter) GetCountry(ctx context.Context, country string) (*service.CountryData, error) {
	info, err := a.Client.GetCountry(ctx, country)
	if err != nil {
		return nil, err
	}

	return &service.CountryData{
		Name: info.Name,
		Code: info.Code,
	}, nil
}

func (a CountryProviderAdapter) GetCities(ctx context.Context, country string) ([]string, error) {
	return a.Client.GetCities(ctx, country)
}

type GeoProviderAdapter struct {
	Client *GeoClient
}

func (a GeoProviderAdapter) GetCoordinates(ctx context.Context, city, countryCode string) (float64, float64, error) {
	return a.Client.GetCoordinates(ctx, city, countryCode)
}
