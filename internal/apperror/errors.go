package apperror

import (
	"errors"
	"fmt"
)

// ErrorCode represents the type of error
type ErrorCode string

const (
	CodeBadRequest      ErrorCode = "BAD_REQUEST"
	CodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	CodeForbidden       ErrorCode = "FORBIDDEN"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	CodeValidationError ErrorCode = "VALIDATION_ERROR"
	CodeUnprocessable   ErrorCode = "UNPROCESSABLE_ENTITY"
)

// AppError is the application-specific error type
type AppError struct {
	Code       ErrorCode
	Message    string
	StatusCode int
	Details    map[string]interface{}
	Err        error // Original error for logging
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Helper functions to create common errors

func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: 400,
		Details:    make(map[string]interface{}),
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: 401,
		Details:    make(map[string]interface{}),
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: 403,
		Details:    make(map[string]interface{}),
	}
}

func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:       CodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: 404,
		Details:    make(map[string]interface{}),
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: 409,
		Details:    make(map[string]interface{}),
	}
}

func NewValidationError(message string, details map[string]string) *AppError {
	detailsMap := make(map[string]interface{})
	for k, v := range details {
		detailsMap[k] = v
	}
	return &AppError{
		Code:       CodeValidationError,
		Message:    message,
		StatusCode: 422,
		Details:    detailsMap,
	}
}

func NewInternalServer(message string, err error) *AppError {
	return &AppError{
		Code:       CodeInternalServer,
		Message:    message,
		StatusCode: 500,
		Details:    make(map[string]interface{}),
		Err:        err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:       CodeDatabaseError,
		Message:    message,
		StatusCode: 500,
		Details:    make(map[string]interface{}),
		Err:        err,
	}
}

// From converts any error to AppError if needed
func From(err error) *AppError {
	if err == nil {
		return nil
	}

	// If already AppError, return as-is
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Otherwise, wrap as internal server error
	return NewInternalServer("An internal server error occurred", err)
}
