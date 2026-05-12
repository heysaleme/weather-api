package service

import (
	"context"
	"sort"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type UserService struct {
	users   UserStore
	cities  CityStore
	history WeatherHistoryStore
}

func NewUserService(
	users UserStore,
	cities CityStore,
	history WeatherHistoryStore,
) *UserService {
	return &UserService{
		users:   users,
		cities:  cities,
		history: history,
	}
}

func (s *UserService) List(ctx context.Context) ([]*model.User, error) {
	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.PasswordHash = ""
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NotFound("user not found")
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) GetCurrent(ctx context.Context, authUser *model.AuthUser) (*model.User, error) {
	return s.GetByID(ctx, authUser.ID)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	user, err := s.users.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errs.NotFound("user not found")
	}

	if err := s.users.Delete(ctx, id); err != nil {
		return err
	}
	if err := s.cities.DeleteByUserID(ctx, id); err != nil {
		return err
	}
	if err := s.history.DeleteByUserID(ctx, id); err != nil {
		return err
	}

	return nil
}
