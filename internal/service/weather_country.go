package service

import (
	"context"
	"errors"
	"sort"
	"sync"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type cityWeatherResult struct {
	weather *model.WeatherResult
	err     error
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
