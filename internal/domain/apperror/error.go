package apperror

import (
	"net/http"
)

type AppError struct {
	Code      int
	ErrorCode string
	Message   string
	Err       error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) WithCode(code string) *AppError {
	e.ErrorCode = code
	return e
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewAppError creates a new AppError with all parameters
func NewAppError(code int, errorCode string, message string, err error) *AppError {
	return &AppError{
		Code:      code,
		ErrorCode: errorCode,
		Message:   message,
		Err:       err,
	}
}

// New creates a basic AppError (for backward compatibility)
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		ErrorCode: UnknownError,
	}
}

// ==================== HTTP Error Helpers ====================

// NewBadRequestError creates a 400 Bad Request error
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:      http.StatusBadRequest,
		ErrorCode: BadRequestError,
		Message:   message,
		Err:       nil,
	}
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:      http.StatusNotFound,
		ErrorCode: NotFoundError,
		Message:   message,
		Err:       nil,
	}
}

// NewUnauthorizedError creates a 401 Unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:      http.StatusUnauthorized,
		ErrorCode: AuthUnauthorized,
		Message:   message,
		Err:       nil,
	}
}

// NewForbiddenError creates a 403 Forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:      http.StatusForbidden,
		ErrorCode: AuthForbidden,
		Message:   message,
		Err:       nil,
	}
}

// NewConflictError creates a 409 Conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:      http.StatusConflict,
		ErrorCode: ConflictError,
		Message:   message,
		Err:       nil,
	}
}

// NewInternalError creates a 500 Internal Server error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:      http.StatusInternalServerError,
		ErrorCode: SystemInternalError,
		Message:   message,
		Err:       err,
	}
}

// NewValidationError creates a validation error
func NewValidationError(err error) *AppError {
	return &AppError{
		Code:      http.StatusBadRequest,
		ErrorCode: ValidationError,
		Message:   "validation failed",
		Err:       err,
	}
}

// ==================== Legacy Helpers (Backward Compatibility) ====================

// BadRequest creates a 400 error (legacy)
func BadRequest(msg string, err error) *AppError {
	return New(http.StatusBadRequest, msg, err).
		WithCode(BadRequestError)
}

// NotFound creates a 404 error (legacy)
func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err).
		WithCode(NotFoundError)
}

// Unauthorized creates a 401 error (legacy)
func Unauthorized(msg string, err error) *AppError {
	return New(http.StatusUnauthorized, msg, err).
		WithCode(AuthUnauthorized)
}

// Internal creates a 500 error (legacy)
func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "internal server error", err).
		WithCode(SystemInternalError)
}

// Validation creates a validation error (legacy)
func Validation(err error) *AppError {
	return New(
		http.StatusBadRequest,
		"validation failed",
		err,
	).WithCode(ValidationError)
}
