package basicauth

import "net/http"

// IAuth is the interface that wraps the Authenticate method.
type IAuth interface {
	Authenticate(req *http.Request) error
}

// SecretProvider is the interface that wraps the GetBasicAuthUsername and GetBasicAuthPassword methods.
type SecretProvider interface {
	GetBasicAuthUsername() []byte
	GetBasicAuthPassword() []byte
}
