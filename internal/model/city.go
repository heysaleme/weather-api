package model

import "time"

type City struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type WeatherHistoryRecord struct {
	ID          int64         `json:"id"`
	UserID      int64         `json:"user_id,omitempty"`
	CityID      int64         `json:"city_id"`
	City        string        `json:"city"`
	Weather     WeatherResult `json:"weather"`
	RequestedAt time.Time     `json:"requested_at"`
}
