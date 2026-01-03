package errors

import "fmt"

// AppError adalah custom error struct untuk aplikasi
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap mengimplementasikan interface errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Err
}

// New membuat AppError baru
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap membungkus error existing dengan AppError
func Wrap(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error codes
const (
	ErrCodeBadRequest          = 400
	ErrCodeUnauthorized        = 401
	ErrCodeForbidden           = 403
	ErrCodeNotFound            = 404
	ErrCodeConflict            = 409
	ErrCodeInternalServerError = 500
)

// Common errors
var (
	ErrBadRequest          = New(ErrCodeBadRequest, "Bad request")
	ErrUnauthorized        = New(ErrCodeUnauthorized, "Unauthorized")
	ErrForbidden           = New(ErrCodeForbidden, "Forbidden")
	ErrNotFound            = New(ErrCodeNotFound, "Resource not found")
	ErrInternalServerError = New(ErrCodeInternalServerError, "Internal server error")
)
