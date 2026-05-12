package service

import (
	"context"
	"sort"
	"strings"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type CityService struct {
	cities CityStore
}

func NewCityService(cities CityStore) *CityService {
	return &CityService{cities: cities}
}

func (s *CityService) Create(ctx context.Context, user *model.AuthUser, cityName string) (*model.City, error) {
	name := strings.TrimSpace(cityName)
	if name == "" {
		return nil, errs.InvalidInput("city is required")
	}

	cities, err := s.cities.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	for _, city := range cities {
		if strings.EqualFold(city.Name, name) {
			return nil, errs.Conflict("city already added")
		}
	}

	return s.cities.Create(ctx, user.ID, name)
}

func (s *CityService) List(ctx context.Context, user *model.AuthUser) ([]*model.City, error) {
	cities, err := s.cities.ListByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	sort.Slice(cities, func(i, j int) bool {
		return cities[i].ID < cities[j].ID
	})

	return cities, nil
}

func (s *CityService) Delete(ctx context.Context, user *model.AuthUser, cityID int64) error {
	if cityID <= 0 {
		return errs.InvalidInput("invalid city id")
	}

	city, err := s.cities.GetByID(ctx, cityID)
	if err != nil {
		return err
	}
	if city == nil || city.UserID != user.ID {
		return errs.NotFound("city not found")
	}

	return s.cities.Delete(ctx, cityID)
}
