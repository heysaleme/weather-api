package service

import (
	"context"

	"weather-api/internal/model"
)

func (s *WeatherService) fetchCityWeather(ctx context.Context, city, countryCode string) (*model.WeatherResult, error) {
	lat, lon, err := s.geoClient.GetCoordinates(ctx, city, countryCode)
	if err != nil {
		return nil, err
	}

	resp, err := s.weatherClient.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, err
	}

	return &model.WeatherResult{
		City:        city,
		Latitude:    lat,
		Longitude:   lon,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
		Clothing:    getClothing(resp.Temperature),
	}, nil
}
