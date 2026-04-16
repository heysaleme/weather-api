package service

import (
	"context"
	"fmt"
	"sort"

	"weather-api/internal/client"
	"weather-api/internal/model"
)

const (
	coldThreshold = 5
	warmThreshold = 15
)

type WeatherService struct {
	weatherClient *client.WeatherClient
	geoClient     *client.GeoClient
}

func NewWeatherService(w *client.WeatherClient, g *client.GeoClient) *WeatherService {
	return &WeatherService{
		weatherClient: w,
		geoClient:     g,
	}
}

func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*model.WeatherResult, error) {
	lat, lon, err := s.geoClient.GetCoordinates(ctx, city)
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
		Temperature: resp.CurrentWeather.Temperature,
		WindSpeed:   resp.CurrentWeather.Windspeed,
		WeatherCode: resp.CurrentWeather.Weathercode,
		Time:        resp.CurrentWeather.Time,
		Description: mapWeatherCode(resp.CurrentWeather.Weathercode),
		Clothing:    getClothing(resp.CurrentWeather.Temperature),
	}, nil
}

func (s *WeatherService) GetWeatherByCountry(ctx context.Context, country string) ([]*model.WeatherResult, error) {
	cities, ok := countryCities[country]
	if !ok {
		return nil, fmt.Errorf("country not supported")
	}

	var results []*model.WeatherResult

	for _, city := range cities {
		w, err := s.GetWeatherByCity(ctx, city)
		if err != nil {
			continue
		}
		results = append(results, w)
	}

	return results, nil
}

func (s *WeatherService) GetTopCitiesByCountry(ctx context.Context, country string) ([]*model.WeatherResult, error) {
	cities, err := s.GetWeatherByCountry(ctx, country)
	if err != nil {
		return nil, err
	}

	sort.Slice(cities, func(i, j int) bool {
		return cities[i].Temperature > cities[j].Temperature
	})

	if len(cities) > 3 {
		return cities[:3], nil
	}

	return cities, nil
}

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	default:
		return "Неизвестно"
	}
}

func getClothing(temp float64) string {
	if temp < coldThreshold {
		return "Тёплая одежда"
	}
	if temp < warmThreshold {
		return "Куртка"
	}
	return "Лёгкая одежда"
}
