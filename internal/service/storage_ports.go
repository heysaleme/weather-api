package service

import (
	"context"

	"weather-api/internal/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type CityStore interface {
	Create(ctx context.Context, userID int64, name string) (*model.City, error)
	ListByUserID(ctx context.Context, userID int64) ([]*model.City, error)
	GetByID(ctx context.Context, id int64) (*model.City, error)
	Delete(ctx context.Context, id int64) error
	DeleteByUserID(ctx context.Context, userID int64) error
}

type WeatherHistoryStore interface {
	Create(ctx context.Context, record *model.WeatherHistoryRecord) error
	ListByUserID(ctx context.Context, userID int64) ([]*model.WeatherHistoryRecord, error)
	DeleteByUserID(ctx context.Context, userID int64) error
}
