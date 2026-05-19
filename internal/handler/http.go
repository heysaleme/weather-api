package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"weather-api/internal/errs"
	"weather-api/internal/middleware"
	"weather-api/internal/model"

	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(buf.Bytes())
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	status := statusCode(err)
	if status >= http.StatusInternalServerError {
		logInternalError(r, err, status)
		writeError(w, status, "internal server error")
		return
	}

	writeError(w, status, err.Error())
}

func readJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func statusCode(err error) int {
	switch {
	case errors.Is(err, errs.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errs.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errs.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errs.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, errs.ErrUpstream):
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

func currentUser(r *http.Request) (*model.AuthUser, error) {
	user, ok := middleware.UserFromContext(r.Context())
	if !ok {
		return nil, errs.Unauthorized("unauthorized")
	}
	return user, nil
}

func logInternalError(r *http.Request, err error, status int) {
	logger, ok := middleware.LoggerFromContext(r.Context())
	if !ok {
		logger = zap.NewNop()
	}

	fields := []zap.Field{
		zap.Error(err),
		zap.Int("status_code", status),
	}
	if requestID, ok := middleware.RequestIDFromContext(r.Context()); ok {
		fields = append(fields, zap.String("request_id", requestID))
	}

	logger.Error("request failed", fields...)
}
