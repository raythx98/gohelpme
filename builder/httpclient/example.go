package httpclient

import (
	"net/http"
)

// Example is an example of a custom RoundTripper that can be used with the http.Client.
type Example struct{}

// RoundTrip is a method that satisfies the http.RoundTripper interface.
// It is called by the http.Client to execute a single HTTP transaction.
// This method can be used to perform custom logic before and after the request is sent.
func (t *Example) RoundTrip(req *http.Request) (*http.Response, error) {
	// Do work before the request is sent

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Do work after the response is received
	return resp, err
}
