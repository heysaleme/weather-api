package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"

	"weather-api/internal/client"
	"weather-api/internal/errs"
	"weather-api/internal/model"
)

const (
	coldThreshold = 5
	warmThreshold = 15
)

type WeatherService struct {
	weatherClient *client.WeatherClient
	geoClient     *client.GeoClient
	countryClient *client.CountryClient
}

func NewWeatherService(w *client.WeatherClient, g *client.GeoClient, c *client.CountryClient) *WeatherService {
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

type cityWeatherResult struct {
	weather *model.WeatherResult
	err     error
}

func (s *WeatherService) loadCountryCities(ctx context.Context, country string) (*client.CountryInfo, []string, error) {
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

func (s *WeatherService) fetchCountryWeather(ctx context.Context, cities []string, countryCode string) ([]*model.WeatherResult, error) {
	out := make(chan cityWeatherResult, len(cities))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go func(city string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			weather, err := s.fetchCityWeather(ctx, city, countryCode)
			out <- cityWeatherResult{weather: weather, err: err}
		}(city)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	results, err := collectCityWeather(out)
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].City < results[j].City
	})

	return results, nil
}

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
		Temperature: resp.CurrentWeather.Temperature,
		WindSpeed:   resp.CurrentWeather.Windspeed,
		WeatherCode: resp.CurrentWeather.Weathercode,
		Time:        resp.CurrentWeather.Time,
		Description: mapWeatherCode(resp.CurrentWeather.Weathercode),
		Clothing:    getClothing(resp.CurrentWeather.Temperature),
	}, nil
}

func collectCityWeather(out <-chan cityWeatherResult) ([]*model.WeatherResult, error) {
	results := make([]*model.WeatherResult, 0)
	var lastErr error

	for item := range out {
		if item.err != nil {
			if !errors.Is(item.err, errs.ErrNotFound) {
				lastErr = item.err
			}
			continue
		}
		results = append(results, item.weather)
	}

	if len(results) > 0 {
		return results, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}

	return nil, errs.NotFound("no weather data found for country")
}

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
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

func uniqueCities(cities []string) []string {
	seen := make(map[string]struct{}, len(cities))
	result := make([]string, 0, len(cities))

	for _, city := range cities {
		city = strings.TrimSpace(city)
		if city == "" {
			continue
		}
		key := strings.ToLower(city)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, city)
	}

	return result
}
