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

type MockCityStore struct {
	mock.Mock
}

func (m *MockCityStore) Create(ctx context.Context, userID int64, name string) (*model.City, error) {
	args := m.Called(ctx, userID, name)
	city, _ := args.Get(0).(*model.City)
	return city, args.Error(1)
}

func (m *MockCityStore) ListByUserID(ctx context.Context, userID int64) ([]*model.City, error) {
	args := m.Called(ctx, userID)
	cities, _ := args.Get(0).([]*model.City)
	return cities, args.Error(1)
}

func (m *MockCityStore) GetByID(ctx context.Context, id int64) (*model.City, error) {
	args := m.Called(ctx, id)
	city, _ := args.Get(0).(*model.City)
	return city, args.Error(1)
}

func (m *MockCityStore) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockCityStore) DeleteByUserID(ctx context.Context, userID int64) error {
	return m.Called(ctx, userID).Error(0)
}

func TestCityServiceCreateSuccess(t *testing.T) {
	t.Parallel()

	store := new(MockCityStore)
	svc := NewCityService(store)
	user := &model.AuthUser{ID: 7}
	expected := &model.City{ID: 11, UserID: 7, Name: "Almaty", CreatedAt: time.Now().UTC()}

	store.On("ListByUserID", mock.Anything, int64(7)).Return([]*model.City{}, nil).Once()
	store.On("Create", mock.Anything, int64(7), "Almaty").Return(expected, nil).Once()

	city, err := svc.Create(context.Background(), user, "  Almaty  ")
	require.NoError(t, err)
	require.NotNil(t, city)
	assert.Equal(t, expected, city)
	store.AssertExpectations(t)
}

func TestCityServiceCreateRejectsEmptyName(t *testing.T) {
	t.Parallel()

	svc := NewCityService(new(MockCityStore))

	city, err := svc.Create(context.Background(), &model.AuthUser{ID: 1}, "   ")
	require.Error(t, err)
	assert.Nil(t, city)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)
}

func TestCityServiceCreateRejectsDuplicateCity(t *testing.T) {
	t.Parallel()

	store := new(MockCityStore)
	svc := NewCityService(store)

	store.On("ListByUserID", mock.Anything, int64(5)).Return([]*model.City{
		{ID: 1, UserID: 5, Name: "almaty"},
	}, nil).Once()

	city, err := svc.Create(context.Background(), &model.AuthUser{ID: 5}, "Almaty")
	require.Error(t, err)
	assert.Nil(t, city)
	assert.ErrorIs(t, err, errs.ErrConflict)
	store.AssertExpectations(t)
}

func TestCityServiceCreatePropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	store := new(MockCityStore)
	svc := NewCityService(store)
	expectedErr := errors.New("list failed")

	store.On("ListByUserID", mock.Anything, int64(5)).Return(nil, expectedErr).Once()

	city, err := svc.Create(context.Background(), &model.AuthUser{ID: 5}, "Almaty")
	require.Error(t, err)
	assert.Nil(t, city)
	assert.ErrorIs(t, err, expectedErr)
	store.AssertExpectations(t)
}

func TestCityServiceDeleteRejectsInvalidID(t *testing.T) {
	t.Parallel()

	svc := NewCityService(new(MockCityStore))

	err := svc.Delete(context.Background(), &model.AuthUser{ID: 7}, 0)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)
}

func TestCityServiceDeleteReturnsNotFoundForForeignCity(t *testing.T) {
	t.Parallel()

	store := new(MockCityStore)
	svc := NewCityService(store)

	store.On("GetByID", mock.Anything, int64(9)).Return(&model.City{ID: 9, UserID: 999}, nil).Once()

	err := svc.Delete(context.Background(), &model.AuthUser{ID: 7}, 9)
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrNotFound)
	store.AssertExpectations(t)
}

func TestCityServiceDeletePropagatesDeleteError(t *testing.T) {
	t.Parallel()

	store := new(MockCityStore)
	svc := NewCityService(store)
	expectedErr := errors.New("delete failed")

	store.On("GetByID", mock.Anything, int64(9)).Return(&model.City{ID: 9, UserID: 7}, nil).Once()
	store.On("Delete", mock.Anything, int64(9)).Return(expectedErr).Once()

	err := svc.Delete(context.Background(), &model.AuthUser{ID: 7}, 9)
	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	store.AssertExpectations(t)
}
