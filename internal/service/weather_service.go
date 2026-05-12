package service

import (
	"context"
	"sort"
	"strings"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

const (
	coldThreshold = 5
	warmThreshold = 15
)

type WeatherService struct {
	weatherClient WeatherProvider
	geoClient     GeoProvider
	countryClient CountryProvider
}

func NewWeatherService(w WeatherProvider, g GeoProvider, c CountryProvider) *WeatherService {
	return &WeatherService{
		weatherClient: w,
		geoClient:     g,
		countryClient: c,
	}
}

func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*model.WeatherResult, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return nil, errs.InvalidInput("city is required")
	}

	return s.fetchCityWeather(ctx, city, "")
}

func (s *WeatherService) GetWeatherByCountry(ctx context.Context, country string) ([]*model.WeatherResult, error) {
	info, cities, err := s.loadCountryCities(ctx, country)
	if err != nil {
		return nil, err
	}

	return s.fetchCountryWeather(ctx, cities, info.Code)
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

func (s *WeatherService) loadCountryCities(ctx context.Context, country string) (*CountryData, []string, error) {
	info, err := s.countryClient.GetCountry(ctx, country)
	if err != nil {
		return nil, nil, err
	}

	cities, err := s.countryClient.GetCities(ctx, info.Name)
	if err != nil {
		return nil, nil, err
	}

	return info, uniqueCities(cities), nil
}
