package repository

import (
	"context"
	"sync"

	"weather-api/internal/model"
)

type WeatherHistoryRepository interface {
	Create(ctx context.Context, record *model.WeatherHistoryRecord) error
	ListByUserID(ctx context.Context, userID int64) ([]*model.WeatherHistoryRecord, error)
	DeleteByUserID(ctx context.Context, userID int64) error
}

type InMemoryWeatherHistoryRepository struct {
	mu     sync.RWMutex
	nextID int64
	items  []*model.WeatherHistoryRecord
}

func NewInMemoryWeatherHistoryRepository() *InMemoryWeatherHistoryRepository {
	return &InMemoryWeatherHistoryRepository{
		nextID: 1,
		items:  make([]*model.WeatherHistoryRecord, 0),
	}
}

func (r *InMemoryWeatherHistoryRepository) Create(_ context.Context, record *model.WeatherHistoryRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	recordCopy := *record
	recordCopy.ID = r.nextID
	r.nextID++
	r.items = append(r.items, &recordCopy)
	record.ID = recordCopy.ID

	return nil
}

func (r *InMemoryWeatherHistoryRepository) ListByUserID(_ context.Context, userID int64) ([]*model.WeatherHistoryRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.WeatherHistoryRecord, 0)
	for _, item := range r.items {
		if item.UserID != userID {
			continue
		}
		recordCopy := *item
		result = append(result, &recordCopy)
	}

	return result, nil
}

func (r *InMemoryWeatherHistoryRepository) DeleteByUserID(_ context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.items[:0]
	for _, item := range r.items {
		if item.UserID != userID {
			filtered = append(filtered, item)
		}
	}
	r.items = filtered
	return nil
}
