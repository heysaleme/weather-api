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
	"weather-api/internal/repository"
)

type stubWeatherProvider struct {
	resp *WeatherData
	err  error
	get  func(lat, lon float64) (*WeatherData, error)
}

func (s stubWeatherProvider) GetCurrentWeather(ctx context.Context, lat, lon float64) (*WeatherData, error) {
	if s.get != nil {
		return s.get(lat, lon)
	}
	return s.resp, s.err
}

type stubGeoProvider struct {
	lat float64
	lon float64
	err error
	get func(city, countryCode string) (float64, float64, error)
}

func (s stubGeoProvider) GetCoordinates(ctx context.Context, city, countryCode string) (float64, float64, error) {
	if s.get != nil {
		return s.get(city, countryCode)
	}
	return s.lat, s.lon, s.err
}

type stubCountryProvider struct {
	info      *CountryData
	cities    []string
	infoErr   error
	citiesErr error
}

func (s stubCountryProvider) GetCountry(ctx context.Context, country string) (*CountryData, error) {
	return s.info, s.infoErr
}

func (s stubCountryProvider) GetCities(ctx context.Context, country string) ([]string, error) {
	return s.cities, s.citiesErr
}

func TestWeatherHelpers(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"Almaty", "Astana"}, uniqueCities([]string{" Almaty ", "almaty", "", "Astana"}))
	assert.Equal(t, "Ясно", mapWeatherCode(0))
	assert.Equal(t, "Неизвестно", mapWeatherCode(999))
	assert.Equal(t, "Тёплая одежда", getClothing(1))
	assert.Equal(t, "Куртка", getClothing(10))
	assert.Equal(t, "Лёгкая одежда", getClothing(20))
}

func TestWeatherServiceGetWeatherByCity(t *testing.T) {
	t.Parallel()

	svc := NewWeatherService(
		stubWeatherProvider{resp: &WeatherData{Temperature: 18, WindSpeed: 4, WeatherCode: 0, Time: "2026-05-19T12:00:00Z"}},
		stubGeoProvider{lat: 43.24, lon: 76.92},
		stubCountryProvider{},
	)

	result, err := svc.GetWeatherByCity(context.Background(), "  Almaty ")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Almaty", result.City)
	assert.Equal(t, 43.24, result.Latitude)
	assert.Equal(t, "Ясно", result.Description)
	assert.Equal(t, "Лёгкая одежда", result.Clothing)
}

func TestWeatherServiceGetWeatherByCityRejectsEmptyName(t *testing.T) {
	t.Parallel()

	svc := NewWeatherService(stubWeatherProvider{}, stubGeoProvider{}, stubCountryProvider{})

	result, err := svc.GetWeatherByCity(context.Background(), "   ")
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)
}

func TestWeatherServiceGetTopCitiesByCountry(t *testing.T) {
	t.Parallel()

	svc := NewWeatherService(
		stubWeatherProvider{
			get: func(lat, lon float64) (*WeatherData, error) {
				switch lat {
				case 1:
					return &WeatherData{Temperature: 21, WindSpeed: 2, WeatherCode: 1, Time: "2026-05-19T12:00:00Z"}, nil
				case 2:
					return &WeatherData{Temperature: 15, WindSpeed: 2, WeatherCode: 1, Time: "2026-05-19T12:00:00Z"}, nil
				default:
					return &WeatherData{Temperature: 8, WindSpeed: 2, WeatherCode: 1, Time: "2026-05-19T12:00:00Z"}, nil
				}
			},
		},
		stubGeoProvider{
			get: func(city, countryCode string) (float64, float64, error) {
				switch city {
				case "Almaty":
					return 1, 1, nil
				case "Astana":
					return 2, 2, nil
				default:
					return 3, 3, nil
				}
			},
		},
		stubCountryProvider{
			info:   &CountryData{Name: "Kazakhstan", Code: "KZ"},
			cities: []string{"Almaty", "Astana", "Shymkent", "Almaty"},
		},
	)

	results, err := svc.GetTopCitiesByCountry(context.Background(), "Kazakhstan")
	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, "Almaty", results[0].City)
	assert.Equal(t, "Astana", results[1].City)
}

func TestCollectCityWeatherHandlesErrors(t *testing.T) {
	t.Parallel()

	results, err := collectCityWeather(chanWithItems(
		cityWeatherResult{err: errs.NotFound("missing")},
		cityWeatherResult{err: errors.New("upstream failed")},
	))
	require.Error(t, err)
	assert.Nil(t, results)
	assert.EqualError(t, err, "upstream failed")

	results, err = collectCityWeather(chanWithItems(
		cityWeatherResult{err: errs.NotFound("missing")},
	))
	require.Error(t, err)
	assert.Nil(t, results)
	assert.ErrorIs(t, err, errs.ErrNotFound)
}

func TestUserWeatherServiceGetCurrentAndHistory(t *testing.T) {
	t.Parallel()

	cities := new(MockCityStore)
	history := new(MockWeatherHistoryStore)
	lookup := stubWeatherLookupService{
		results: map[string]*model.WeatherResult{
			"Astana": {City: "Astana", Temperature: 8},
			"Almaty": {City: "Almaty", Temperature: 18},
		},
	}

	svc := NewUserWeatherService(cities, history, lookup)
	user := &model.AuthUser{ID: 3}

	cities.On("ListByUserID", mock.Anything, int64(3)).Return([]*model.City{
		{ID: 1, UserID: 3, Name: "Astana"},
		{ID: 2, UserID: 3, Name: "Almaty"},
	}, nil).Once()
	history.On("Create", mock.Anything, recordWithCity("Astana")).Return(nil).Once()
	history.On("Create", mock.Anything, recordWithCity("Almaty")).Return(nil).Once()

	current, err := svc.GetCurrent(context.Background(), user)
	require.NoError(t, err)
	require.Len(t, current, 2)
	assert.Equal(t, "Almaty", current[0].City)
	assert.Equal(t, "Astana", current[1].City)

	now := time.Now().UTC()
	history.On("ListByUserID", mock.Anything, int64(3)).Return([]*model.WeatherHistoryRecord{
		{City: "Old", RequestedAt: now.Add(-time.Hour)},
		{City: "New", RequestedAt: now},
	}, nil).Once()

	records, err := svc.GetHistory(context.Background(), user)
	require.NoError(t, err)
	require.Len(t, records, 2)
	assert.Equal(t, "New", records[0].City)
	cities.AssertExpectations(t)
	history.AssertExpectations(t)
}

func TestUserWeatherServiceGetCurrentErrors(t *testing.T) {
	t.Parallel()

	cities := new(MockCityStore)
	history := new(MockWeatherHistoryStore)
	lookup := stubWeatherLookupService{}
	svc := NewUserWeatherService(cities, history, lookup)
	user := &model.AuthUser{ID: 3}

	cities.On("ListByUserID", mock.Anything, int64(3)).Return([]*model.City{}, nil).Once()

	current, err := svc.GetCurrent(context.Background(), user)
	require.Error(t, err)
	assert.Nil(t, current)
	assert.ErrorIs(t, err, errs.ErrNotFound)

	cities2 := new(MockCityStore)
	history2 := new(MockWeatherHistoryStore)
	lookup2 := stubWeatherLookupService{err: errors.New("lookup failed")}
	svc2 := NewUserWeatherService(cities2, history2, lookup2)
	cities2.On("ListByUserID", mock.Anything, int64(3)).Return([]*model.City{
		{ID: 1, UserID: 3, Name: "Astana"},
	}, nil).Once()

	current, err = svc2.GetCurrent(context.Background(), user)
	require.Error(t, err)
	assert.Nil(t, current)
	assert.EqualError(t, err, "lookup failed")
}

func TestAuthServiceEnsureAdminAccountValidationAndReuse(t *testing.T) {
	t.Parallel()

	repo := repository.NewInMemoryUserRepository()
	svc := NewAuthService(repo, fakeTokenManager{})

	require.NoError(t, svc.EnsureAdminAccount(context.Background(), "", ""))

	err := svc.EnsureAdminAccount(context.Background(), "admin@example.com", "short")
	require.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrInvalidInput)

	require.NoError(t, svc.EnsureAdminAccount(context.Background(), "admin@example.com", "very-strong-pass"))
	require.NoError(t, svc.EnsureAdminAccount(context.Background(), "admin@example.com", "very-strong-pass"))
}

type stubWeatherLookupService struct {
	results map[string]*model.WeatherResult
	err     error
}

func (s stubWeatherLookupService) GetWeatherByCity(ctx context.Context, city string) (*model.WeatherResult, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.results[city], nil
}

func chanWithItems(items ...cityWeatherResult) <-chan cityWeatherResult {
	ch := make(chan cityWeatherResult, len(items))
	for _, item := range items {
		ch <- item
	}
	close(ch)
	return ch
}

func recordWithCity(city string) any {
	return mock.MatchedBy(func(record *model.WeatherHistoryRecord) bool {
		return record != nil && record.City == city
	})
}
