package service

import (
	"context"
	"sort"

	"weather-api/internal/model"
)

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
