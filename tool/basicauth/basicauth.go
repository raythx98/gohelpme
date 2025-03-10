package basicauth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Auth is a basic auth implementation.
type Auth struct {
	expectedToken string
}

// New creates a new Auth instance.
func New(configProvider ConfigProvider) *Auth {
	return &Auth{
		expectedToken: base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", configProvider.GetBasicAuthUsername(), configProvider.GetBasicAuthPassword()))),
	}
}

// Authenticate authenticates the request using basic auth.
//
// It returns an error if the request is not authenticated.
func (a *Auth) Authenticate(req *http.Request) error {
	tokenHeader := req.Header.Get("Authorization")

	// The usual convention is for "Bearer" to be title-cased. However, there's no
	// strict rule around this, and it's best to follow the robustness principle here.
	if len(tokenHeader) < 6 || !strings.EqualFold(tokenHeader[:6], "basic ") {
		return fmt.Errorf("no token present in request")
	}

	if a.expectedToken != tokenHeader[6:] {
		return fmt.Errorf("invalid token")
	}

	return nil
}
