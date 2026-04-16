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

	writeJSON(w, http.StatusOK, result)
}

func (h *WeatherHandler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetWeatherByCountry(r.Context(), country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *WeatherHandler) GetTopCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")

	result, err := h.service.GetTopCitiesByCountry(r.Context(), country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
	}
}
