package errorhelper

import "fmt"

// AuthError is an error type for unauthorized access
type AuthError struct {
	Err error
}

// NewAuthError creates a new AuthError
func NewAuthError(err error) *AuthError {
	return &AuthError{
		Err: err,
	}
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("Unauthorized, Err: %v", e.Err)
}
