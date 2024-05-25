package httpclient

import (
	"net/http"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	// TODO: should this accept an ILogger interface?
	httpClient := &http.Client{
		Transport: &LogRoundTripper{},
	}
	return &Client{httpClient: httpClient}
}

func (i *Client) Do(req *http.Request) (*http.Response, error) {
	return i.httpClient.Do(req)
}
