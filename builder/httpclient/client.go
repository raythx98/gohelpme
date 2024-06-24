package httpclient

import (
	"github.com/raythx98/gohelpme/tool/logger"
	"net/http"
)

// Client is a wrapper around http.Client that logs requests and responses.
type Client struct {
	httpClient *http.Client // TODO: Create a IHttpClient interface
}

// New creates a new Client.
func New(log logger.ILogger) *Client {
	httpClient := &http.Client{
		Transport: NewLogRoundTripper(log),
	}
	return &Client{httpClient: httpClient}
}

// Do sends an HTTP request and returns an HTTP response, following policy (like redirects, cookies, auth) as configured on the client.
func (i *Client) Do(req *http.Request) (*http.Response, error) {
	return i.httpClient.Do(req)
}
