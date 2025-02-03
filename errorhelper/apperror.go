package errorhelper

import "fmt"

// AppError are analogous to caught/known exceptions
type AppError struct {
	code    int
	message string
	err     error
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d, Err: %v", e.message, e.code, e.err)
}
