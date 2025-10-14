package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-level error
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string, statusCode int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// Common error constructors

// ErrBadRequest creates a 400 Bad Request error
func ErrBadRequest(message string, err error) *AppError {
	return New("bad_request", message, http.StatusBadRequest, err)
}

// ErrUnauthorized creates a 401 Unauthorized error
func ErrUnauthorized(message string, err error) *AppError {
	return New("unauthorized", message, http.StatusUnauthorized, err)
}

// ErrForbidden creates a 403 Forbidden error
func ErrForbidden(message string, err error) *AppError {
	return New("forbidden", message, http.StatusForbidden, err)
}

// ErrNotFound creates a 404 Not Found error
func ErrNotFound(message string, err error) *AppError {
	return New("not_found", message, http.StatusNotFound, err)
}

// ErrConflict creates a 409 Conflict error
func ErrConflict(message string, err error) *AppError {
	return New("conflict", message, http.StatusConflict, err)
}

// ErrInternal creates a 500 Internal Server Error
func ErrInternal(message string, err error) *AppError {
	return New("internal_error", message, http.StatusInternalServerError, err)
}

// ErrServiceUnavailable creates a 503 Service Unavailable error
func ErrServiceUnavailable(message string, err error) *AppError {
	return New("service_unavailable", message, http.StatusServiceUnavailable, err)
}

// Predefined errors
var (
	ErrInvalidInput      = ErrBadRequest("Invalid input", nil)
	ErrMissingAuthHeader = ErrUnauthorized("Missing authorization header", nil)
	ErrInvalidToken      = ErrUnauthorized("Invalid token", nil)
	ErrExpiredToken      = ErrUnauthorized("Token expired", nil)
	ErrResourceNotFound  = ErrNotFound("Resource not found", nil)
	ErrDatabaseError     = ErrInternal("Database error", nil)
	ErrStorageError      = ErrInternal("Storage error", nil)
)
