package handler

import (
	"context"
	"net/http"

	"weather-api/internal/dto"
	"weather-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type WeatherService interface {
	GetWeatherByCity(ctx context.Context, city string) (*model.WeatherResult, error)
	GetWeatherByCountry(ctx context.Context, country string) ([]*model.WeatherResult, error)
	GetTopCitiesByCountry(ctx context.Context, country string) ([]*model.WeatherResult, error)
}

type WeatherHandler struct {
	service WeatherService
}

func NewWeatherHandler(s WeatherService) *WeatherHandler {
	return &WeatherHandler{service: s}
}

func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	result, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToWeatherResponse(result))
}

func (h *WeatherHandler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetWeatherByCountry(r.Context(), country)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToWeatherResponses(result))
}

func (h *WeatherHandler) GetTopCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetTopCitiesByCountry(r.Context(), country)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToWeatherResponses(result))
}
