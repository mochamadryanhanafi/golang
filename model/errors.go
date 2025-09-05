package model

import "fmt"

// AppError adalah tipe error kustom untuk aplikasi ini.
type AppError struct {
	StatusCode int
	Message    string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// NewAppError membuat instance AppError baru.
func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Pre-defined errors
var (
	ErrInvalidCredentials = NewAppError(401, "invalid email or password")
	ErrAccountNotVerified = NewAppError(403, "account is not verified")
	ErrUserAlreadyExists  = NewAppError(409, "user with this email already exists")
	ErrUserNotFound       = NewAppError(404, "user not found")
	ErrInvalidOTP         = NewAppError(400, "invalid or expired OTP")
	ErrInvalidToken       = NewAppError(401, "invalid or expired token")
)
