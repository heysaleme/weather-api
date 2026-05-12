package dto

import "time"

type WeatherResponse struct {
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
	Clothing    string  `json:"clothing"`
}

type WeatherHistoryResponse struct {
	ID          int64           `json:"id"`
	CityID      int64           `json:"city_id"`
	City        string          `json:"city"`
	Weather     WeatherResponse `json:"weather"`
	RequestedAt time.Time       `json:"requested_at"`
}
