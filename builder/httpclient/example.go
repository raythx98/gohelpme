package httpclient

import (
	"net/http"
)

type Example struct{}

func (t Example) RoundTrip(req *http.Request) (*http.Response, error) {
	// Do work before the request is sent

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Do work after the response is received
	return resp, err
}
