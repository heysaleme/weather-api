package service

type UserWeatherService struct {
	cities         CityStore
	history        WeatherHistoryStore
	weatherService WeatherLookupService
}

func NewUserWeatherService(
	cities CityStore,
	history WeatherHistoryStore,
	weatherService WeatherLookupService,
) *UserWeatherService {
	return &UserWeatherService{
		cities:         cities,
		history:        history,
		weatherService: weatherService,
	}
}
