// Package apperr defines application errors that translate into consistent
// JSON responses with Russian messages, mirroring the original Flask backend.
package apperr

import "net/http"

// AppError is an expected application error carrying an HTTP status code.
type AppError struct {
	Message    string
	StatusCode int
}

func (e *AppError) Error() string { return e.Message }

// New builds an AppError with an explicit status code.
func New(message string, status int) *AppError {
	return &AppError{Message: message, StatusCode: status}
}

// NotFound -> 404
func NotFound(message string) *AppError { return New(message, http.StatusNotFound) }

// Validation -> 422
func Validation(message string) *AppError { return New(message, http.StatusUnprocessableEntity) }

// Auth -> 401
func Auth(message string) *AppError { return New(message, http.StatusUnauthorized) }

// Permission -> 403
func Permission(message string) *AppError { return New(message, http.StatusForbidden) }

// Conflict -> 409
func Conflict(message string) *AppError { return New(message, http.StatusConflict) }

// BusinessRule -> 422
func BusinessRule(message string) *AppError { return New(message, http.StatusUnprocessableEntity) }

// Bad -> 400
func Bad(message string) *AppError { return New(message, http.StatusBadRequest) }
