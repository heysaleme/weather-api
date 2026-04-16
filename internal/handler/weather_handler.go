package handler

import (
	"context"
	"encoding/json"
	"net/http"

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

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	result, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *WeatherHandler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetWeatherByCountry(r.Context(), country)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *WeatherHandler) GetTopCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetTopCitiesByCountry(r.Context(), country)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to encode json")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
