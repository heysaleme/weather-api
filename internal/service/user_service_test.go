package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"weather-api/internal/errs"
	"weather-api/internal/model"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Create(ctx context.Context, user *model.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *MockUserStore) GetByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*model.User)
	return user, args.Error(1)
}

func (m *MockUserStore) List(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	users, _ := args.Get(0).([]*model.User)
	return users, args.Error(1)
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type MockWeatherHistoryStore struct {
	mock.Mock
}

func (m *MockWeatherHistoryStore) Create(ctx context.Context, record *model.WeatherHistoryRecord) error {
	return m.Called(ctx, record).Error(0)
}

func (m *MockWeatherHistoryStore) ListByUserID(ctx context.Context, userID int64) ([]*model.WeatherHistoryRecord, error) {
	args := m.Called(ctx, userID)
	records, _ := args.Get(0).([]*model.WeatherHistoryRecord)
	return records, args.Error(1)
}

func (m *MockWeatherHistoryStore) DeleteByUserID(ctx context.Context, userID int64) error {
	return m.Called(ctx, userID).Error(0)
}

func TestUserServiceListSanitizesAndSortsUsers(t *testing.T) {
	t.Parallel()

	users := new(MockUserStore)
	cities := new(MockCityStore)
	history := new(MockWeatherHistoryStore)
	svc := NewUserService(users, cities, history)

	users.On("List", mock.Anything).Return([]*model.User{
		{ID: 5, Email: "b@example.com", PasswordHash: "hash-b", CreatedAt: time.Now().UTC()},
		{ID: 2, Email: "a@example.com", PasswordHash: "hash-a", CreatedAt: time.Now().UTC()},
	}, nil).Once()

	result, err := svc.List(context.Background())
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, int64(2), result[0].ID)
	assert.Empty(t, result[0].PasswordHash)
	assert.Empty(t, result[1].PasswordHash)
	users.AssertExpectations(t)
}

func TestUserServiceGetByIDReturnsNotFound(t *testing.T) {
	t.Parallel()

	users := new(MockUserStore)
	svc := NewUserService(users, new(MockCityStore), new(MockWeatherHistoryStore))

	users.On("GetByID", mock.Anything, int64(22)).Return(nil, nil).Once()

	user, err := svc.GetByID(context.Background(), 22)
	require.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errs.ErrNotFound)
	users.AssertExpectations(t)
}

func TestUserServiceDeletePropagatesHistoryError(t *testing.T) {
	t.Parallel()

	users := new(MockUserStore)
	cities := new(MockCityStore)
	history := new(MockWeatherHistoryStore)
	svc := NewUserService(users, cities, history)
	expectedErr := errors.New("history delete failed")

	users.On("GetByID", mock.Anything, int64(3)).Return(&model.User{ID: 3}, nil).Once()
	users.On("Delete", mock.Anything, int64(3)).Return(nil).Once()
	cities.On("DeleteByUserID", mock.Anything, int64(3)).Return(nil).Once()
	history.On("DeleteByUserID", mock.Anything, int64(3)).Return(expectedErr).Once()

	err := svc.Delete(context.Background(), 3)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	users.AssertExpectations(t)
	cities.AssertExpectations(t)
	history.AssertExpectations(t)
}
