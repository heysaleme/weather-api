package handler

import (
	"context"
	"net/http"

	"weather-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type CityService interface {
	Create(ctx context.Context, user *model.AuthUser, req model.CreateCityRequest) (*model.City, error)
	List(ctx context.Context, user *model.AuthUser) ([]*model.City, error)
	Delete(ctx context.Context, user *model.AuthUser, cityID string) error
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

	var req model.CreateCityRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	city, err := h.service.Create(r.Context(), user, req)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, city)
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

	writeJSON(w, http.StatusOK, cities)
}

func (h *CityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, err := currentUser(r)
	if err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), user, chi.URLParam(r, "city_id")); err != nil {
		writeError(w, statusCode(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
