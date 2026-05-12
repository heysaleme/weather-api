package dto

import "time"

type CreateCityRequest struct {
	City string `json:"city"`
}

type CityResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
