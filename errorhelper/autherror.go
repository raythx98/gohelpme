package errorhelper

import "fmt"

// AuthError is an error type for unauthorized access
type AuthError struct {
	Err error
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Unauthorized, Err: %v", e.Err)
}
