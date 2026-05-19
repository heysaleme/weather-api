package handler

import (
	"context"
	"net/http"

	"weather-api/internal/dto"
	"weather-api/internal/model"
)

type UserWeatherService interface {
	GetCurrent(ctx context.Context, user *model.AuthUser) ([]*model.WeatherResult, error)
	GetHistory(ctx context.Context, user *model.AuthUser) ([]*model.WeatherHistoryRecord, error)
}

type UserWeatherHandler struct {
	service UserWeatherService
}

func NewUserWeatherHandler(service UserWeatherService) *UserWeatherHandler {
	return &UserWeatherHandler{service: service}
}

func (h *UserWeatherHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	result, err := h.service.GetCurrent(r.Context(), user)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToWeatherResponses(result))
}

func (h *UserWeatherHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	result, err := h.service.GetHistory(r.Context(), user)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToWeatherHistoryResponses(result))
}
