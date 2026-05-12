package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type weatherItem struct {
	result *model.WeatherResult
	record *model.WeatherHistoryRecord
	err    error
}

func (s *UserWeatherService) GetCurrent(ctx context.Context, user *model.AuthUser) ([]*model.WeatherResult, error) {
	cities, err := s.cities.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, errs.NotFound("no cities found for user")
	}

	items := s.fetchWeatherItems(ctx, user, cities)
	return s.collectCurrentResults(ctx, items, len(cities))
}

func (s *UserWeatherService) fetchWeatherItems(ctx context.Context, user *model.AuthUser, cities []*model.City) <-chan weatherItem {
	out := make(chan weatherItem, len(cities))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go func(city *model.City) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result, err := s.weatherService.GetWeatherByCity(ctx, city.Name)
			if err != nil {
				out <- weatherItem{err: err}
				return
			}

			out <- weatherItem{
				result: result,
				record: &model.WeatherHistoryRecord{
					UserID:      user.ID,
					CityID:      city.ID,
					City:        city.Name,
					Weather:     *result,
					RequestedAt: time.Now().UTC(),
				},
			}
		}(city)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *UserWeatherService) collectCurrentResults(ctx context.Context, items <-chan weatherItem, capacity int) ([]*model.WeatherResult, error) {
	results := make([]*model.WeatherResult, 0, capacity)
	var lastErr error

	for item := range items {
		if item.err != nil {
			lastErr = item.err
			continue
		}
		if err := s.history.Create(ctx, item.record); err != nil {
			return nil, err
		}
		results = append(results, item.result)
	}

	if len(results) == 0 && lastErr != nil {
		return nil, lastErr
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].City < results[j].City
	})

	return results, nil
}
