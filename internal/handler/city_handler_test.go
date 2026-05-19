package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"weather-api/internal/errs"
	"weather-api/internal/middleware"
	"weather-api/internal/model"
)

type MockCityService struct {
	mock.Mock
}

func (m *MockCityService) Create(ctx context.Context, user *model.AuthUser, cityName string) (*model.City, error) {
	args := m.Called(ctx, user, cityName)
	city, _ := args.Get(0).(*model.City)
	return city, args.Error(1)
}

func (m *MockCityService) List(ctx context.Context, user *model.AuthUser) ([]*model.City, error) {
	args := m.Called(ctx, user)
	cities, _ := args.Get(0).([]*model.City)
	return cities, args.Error(1)
}

func (m *MockCityService) Delete(ctx context.Context, user *model.AuthUser, cityID int64) error {
	return m.Called(ctx, user, cityID).Error(0)
}

func TestCityHandlerCreateSuccess(t *testing.T) {
	t.Parallel()

	service := new(MockCityService)
	handler := NewCityHandler(service)
	user := &model.AuthUser{ID: 12, Email: "user@example.com", Role: model.RoleUser}
	createdAt := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)

	service.On("Create", mock.Anything, user, "Almaty").Return(&model.City{
		ID:        4,
		UserID:    12,
		Name:      "Almaty",
		CreatedAt: createdAt,
	}, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/cities", bytes.NewBufferString(`{"city":"Almaty"}`))
	req = req.WithContext(middleware.WithUser(req.Context(), user))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/cities", handler.Create)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&body))
	assert.Equal(t, float64(4), body["id"])
	assert.Equal(t, "Almaty", body["name"])
	service.AssertExpectations(t)
}

func TestCityHandlerCreateRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	handler := NewCityHandler(new(MockCityService))
	user := &model.AuthUser{ID: 12}

	req := httptest.NewRequest(http.MethodPost, "/cities", bytes.NewBufferString(`{"city":`))
	req = req.WithContext(middleware.WithUser(req.Context(), user))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Post("/cities", handler.Create)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error":"invalid request body"}`, rr.Body.String())
}

func TestCityHandlerListSuccess(t *testing.T) {
	t.Parallel()

	service := new(MockCityService)
	handler := NewCityHandler(service)
	user := &model.AuthUser{ID: 12, Email: "user@example.com", Role: model.RoleUser}
	createdAt := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)

	service.On("List", mock.Anything, user).Return([]*model.City{
		{ID: 1, UserID: 12, Name: "Almaty", CreatedAt: createdAt},
	}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	req = req.WithContext(middleware.WithUser(req.Context(), user))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/cities", handler.List)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"id":1,"name":"Almaty","created_at":"2026-05-19T12:00:00Z"}]`, rr.Body.String())
	service.AssertExpectations(t)
}

func TestCityHandlerListReturnsInternalError(t *testing.T) {
	t.Parallel()

	service := new(MockCityService)
	handler := NewCityHandler(service)
	user := &model.AuthUser{ID: 12}

	service.On("List", mock.Anything, user).Return(nil, errors.New("database down")).Once()

	req := httptest.NewRequest(http.MethodGet, "/cities", nil)
	req = req.WithContext(middleware.WithLogger(middleware.WithUser(req.Context(), user), zap.NewNop()))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Get("/cities", handler.List)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.JSONEq(t, `{"error":"internal server error"}`, rr.Body.String())
	service.AssertExpectations(t)
}

func TestCityHandlerDeleteRejectsBadID(t *testing.T) {
	t.Parallel()

	handler := NewCityHandler(new(MockCityService))
	user := &model.AuthUser{ID: 12}

	req := httptest.NewRequest(http.MethodDelete, "/cities/not-a-number", nil)
	req = req.WithContext(middleware.WithUser(req.Context(), user))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/cities/{city_id}", handler.Delete)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, `{"error":"invalid city id"}`, rr.Body.String())
}

func TestCityHandlerDeleteReturnsNotFound(t *testing.T) {
	t.Parallel()

	service := new(MockCityService)
	handler := NewCityHandler(service)
	user := &model.AuthUser{ID: 12}

	service.On("Delete", mock.Anything, user, int64(9)).Return(errs.NotFound("city not found")).Once()

	req := httptest.NewRequest(http.MethodDelete, "/cities/9", nil)
	req = req.WithContext(middleware.WithUser(req.Context(), user))
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Delete("/cities/{city_id}", handler.Delete)
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	assert.JSONEq(t, `{"error":"not found: city not found"}`, rr.Body.String())
	service.AssertExpectations(t)
}
