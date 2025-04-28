package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	ErrCodeBadRequest    = 400
	ErrCodeUnauthorized = 401
	ErrCodeForbidden    = 403
	ErrCodeNotFound     = 404
	ErrCodeConflict     = 409
	ErrCodeValidation   = 422

	ErrCodeInternalServer       = 500
	ErrCodeServiceUnavailable   = 503
	ErrCodeDatabaseError        = 510
	ErrCodeExternalServiceError = 511
)

type AppError struct {
	Code    int         `json:"-"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	Err     error       `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: ", e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func GetStatusCode(err error) int {
	if appErr, ok := AsAppError(err); ok {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

func NewBadRequest(err error, message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Err:     err,
	}
}

func NewBadRequestWithDetails(err error, message string, details interface{}) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func NewNotFound(err error, message string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: message,
		Err:     err,
	}
}

func NewConflict(err error, message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(err error, details interface{}) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: "Validation failed",
		Details: details,
		Err:     err,
	}
}

func NewUnauthorized(err error, message string) *AppError {
	if message == "" {
		message = "Unauthorized access"
	}
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Err:     err,
	}
}

func NewForbidden(err error, message string) *AppError {
	if message == "" {
		message = "Forbidden access"
	}
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
		Err:     err,
	}
}

func NewInternalServerError(err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternalServer,
		Message: "Internal server error",
		Err:     err,
	}
}

func NewDatabaseError(err error, message string) *AppError {
	if message == "" {
		message = "Database operation failed"
	}
	return &AppError{
		Code:    ErrCodeDatabaseError,
		Message: message,
		Err:     err,
	}
}