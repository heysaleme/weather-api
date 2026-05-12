package handler

import (
	"context"
	"net/http"
	"strconv"

	"weather-api/internal/dto"
	"weather-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type CityService interface {
	Create(ctx context.Context, user *model.AuthUser, cityName string) (*model.City, error)
	List(ctx context.Context, user *model.AuthUser) ([]*model.City, error)
	Delete(ctx context.Context, user *model.AuthUser, cityID int64) error
}

type CityHandler struct {
	service CityService
}

func NewCityHandler(service CityService) *CityHandler {
	return &CityHandler{service: service}
}

func (h *CityHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	var req dto.CreateCityRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	city, err := h.service.Create(r.Context(), user, req.City)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, dto.ToCityResponse(city))
}

func (h *CityHandler) List(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	cities, err := h.service.List(r.Context(), user)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.ToCityResponses(cities))
}

func (h *CityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	cityID, err := strconv.ParseInt(chi.URLParam(r, "city_id"), 10, 64)
	if err != nil || cityID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid city id")
		return
	}

	if err := h.service.Delete(r.Context(), user, cityID); err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
