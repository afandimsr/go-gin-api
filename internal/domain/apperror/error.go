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
	// if e.Err != nil {
	// 	return e.Err.Error()
	// }
	return e.Message
}

func (e *AppError) WithCode(code string) *AppError {
	e.ErrorCode = code
	return e
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		ErrorCode: UnknownError,
	}
}

// Helpers
func BadRequest(msg string, err error) *AppError {
	return New(http.StatusBadRequest, msg, err).
		WithCode("BAD_REQUEST")
}

func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err).
		WithCode("NOT_FOUND")
}

func Unauthorized(msg string, err error) *AppError {
	return New(http.StatusUnauthorized, msg, err).
		WithCode(AuthUnauthorized)
}

func Internal(err error) *AppError {
	return New(http.StatusInternalServerError, "internal server error", err).
		WithCode(SystemInternalError)
}

func Validation(err error) *AppError {
	return New(
		http.StatusBadRequest,
		"validasi gagal",
		err,
	)
}
