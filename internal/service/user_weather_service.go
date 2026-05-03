package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"weather-api/internal/errs"
	"weather-api/internal/model"
	"weather-api/internal/repository"
)

type UserWeatherService struct {
	cities         repository.CityRepository
	history        repository.WeatherHistoryRepository
	weatherService *WeatherService
}

func NewUserWeatherService(
	cities repository.CityRepository,
	history repository.WeatherHistoryRepository,
	weatherService *WeatherService,
) *UserWeatherService {
	return &UserWeatherService{
		cities:         cities,
		history:        history,
		weatherService: weatherService,
	}
}

func (s *UserWeatherService) GetCurrent(ctx context.Context, user *model.AuthUser) ([]*model.WeatherResult, error) {
	cities, err := s.cities.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, errs.NotFound("no cities found for user")
	}

	type item struct {
		result *model.WeatherResult
		record *model.WeatherHistoryRecord
		err    error
	}

	out := make(chan item, len(cities))
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
				out <- item{err: err}
				return
			}

			out <- item{
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

	results := make([]*model.WeatherResult, 0, len(cities))
	var lastErr error
	for item := range out {
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

func (s *UserWeatherService) GetHistory(ctx context.Context, user *model.AuthUser) ([]*model.WeatherHistoryRecord, error) {
	records, err := s.history.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].RequestedAt.After(records[j].RequestedAt)
	})

	return records, nil
}
