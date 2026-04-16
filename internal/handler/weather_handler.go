package handler

import (
	"encoding/json"
	"net/http"

	"weather-api/internal/service"

	"github.com/go-chi/chi/v5"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(s *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: s}
}

func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")

	result, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, result)
}

func (h *WeatherHandler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, _ := h.service.GetWeatherByCountry(r.Context(), country)
	writeJSON(w, result)
}

func (h *WeatherHandler) GetTopCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, _ := h.service.GetTopCitiesByCountry(r.Context(), country)
	writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
