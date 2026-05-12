package dto

import "weather-api/internal/model"

func ToUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func ToUserResponses(users []*model.User) []*UserResponse {
	result := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, ToUserResponse(user))
	}
	return result
}

func ToCityResponse(city *model.City) *CityResponse {
	if city == nil {
		return nil
	}

	return &CityResponse{
		ID:        city.ID,
		Name:      city.Name,
		CreatedAt: city.CreatedAt,
	}
}

func ToCityResponses(cities []*model.City) []*CityResponse {
	result := make([]*CityResponse, 0, len(cities))
	for _, city := range cities {
		result = append(result, ToCityResponse(city))
	}
	return result
}

func ToAuthResponse(token string) *AuthResponse {
	return &AuthResponse{AccessToken: token}
}

func ToWeatherResponse(weather *model.WeatherResult) *WeatherResponse {
	if weather == nil {
		return nil
	}

	return &WeatherResponse{
		City:        weather.City,
		Latitude:    weather.Latitude,
		Longitude:   weather.Longitude,
		Temperature: weather.Temperature,
		WindSpeed:   weather.WindSpeed,
		WeatherCode: weather.WeatherCode,
		Time:        weather.Time,
		Description: weather.Description,
		Clothing:    weather.Clothing,
	}
}

func ToWeatherResponses(items []*model.WeatherResult) []*WeatherResponse {
	result := make([]*WeatherResponse, 0, len(items))
	for _, item := range items {
		result = append(result, ToWeatherResponse(item))
	}
	return result
}

func ToWeatherHistoryResponse(record *model.WeatherHistoryRecord) *WeatherHistoryResponse {
	if record == nil {
		return nil
	}

	return &WeatherHistoryResponse{
		ID:          record.ID,
		CityID:      record.CityID,
		City:        record.City,
		Weather:     *ToWeatherResponse(&record.Weather),
		RequestedAt: record.RequestedAt,
	}
}

func ToWeatherHistoryResponses(records []*model.WeatherHistoryRecord) []*WeatherHistoryResponse {
	result := make([]*WeatherHistoryResponse, 0, len(records))
	for _, record := range records {
		result = append(result, ToWeatherHistoryResponse(record))
	}
	return result
}
