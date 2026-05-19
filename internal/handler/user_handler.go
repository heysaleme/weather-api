package handler

import (
	"context"
	"net/http"
	"strconv"

	"weather-api/internal/dto"
	"weather-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type UserService interface {
	List(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetCurrent(ctx context.Context, authUser *model.AuthUser) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToUserResponses(users))
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		respondWithError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToUserResponse(user))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		respondWithError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, http.StatusOK, dto.ToUserResponse(result))
}
