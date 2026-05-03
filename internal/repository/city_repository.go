package repository

import (
	"context"
	"sync"
	"time"

	"weather-api/internal/model"
)

type CityRepository interface {
	Create(ctx context.Context, userID int64, name string) (*model.City, error)
	ListByUserID(ctx context.Context, userID int64) ([]*model.City, error)
	GetByID(ctx context.Context, id int64) (*model.City, error)
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type InMemoryCityRepository struct {
	mu     sync.RWMutex
	nextID int64
	byID   map[int64]*model.City
}

func NewInMemoryCityRepository() *InMemoryCityRepository {
	return &InMemoryCityRepository{
		nextID: 1,
		byID:   make(map[int64]*model.City),
	}
}

func (r *InMemoryCityRepository) Create(_ context.Context, userID int64, name string) (*model.City, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	city := &model.City{
		ID:        r.nextID,
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	r.nextID++
	r.byID[city.ID] = city

	cityCopy := *city
	return &cityCopy, nil
}

func (r *InMemoryCityRepository) ListByUserID(_ context.Context, userID int64) ([]*model.City, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*model.City, 0)
	for _, city := range r.byID {
		if city.UserID != userID {
			continue
		}
		cityCopy := *city
		result = append(result, &cityCopy)
	}

	return result, nil
}

func (r *InMemoryCityRepository) GetByID(_ context.Context, id int64) (*model.City, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	city, ok := r.byID[id]
	if !ok {
		return nil, nil
	}

	cityCopy := *city
	return &cityCopy, nil
}

func (r *InMemoryCityRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.byID, id)
	return nil
}

func (r *InMemoryCityRepository) DeleteByUserID(_ context.Context, userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, city := range r.byID {
		if city.UserID == userID {
			delete(r.byID, id)
		}
	}

	return nil
}
