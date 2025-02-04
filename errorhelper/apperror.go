package errorhelper

import "fmt"

// AppError are analogous to caught/known exceptions
type AppError struct {
	Code    int
	Message string
	err     error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d, Err: %v", e.Message, e.Code, e.err)
}
