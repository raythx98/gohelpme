package errorhelper

import (
	"fmt"
)

// AppError are analogous to caught/known exceptions
type AppError struct {
	Code    int
	Message string
	err     error
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("App Error, Message: %s, Code: %d, Err: %v", e.Message, e.Code, e.err)
}
