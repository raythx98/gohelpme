package httprequest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Builder struct {
	HttpRequest *http.Request
	Error       error
}

// New creates a new http request with the given method and url.
func (i *Builder) New(ctx context.Context, method Method, url string) *Builder {
	req, err := http.NewRequestWithContext(ctx, string(method), url, nil)
	return &Builder{HttpRequest: req, Error: err}
}

// WithBody sets the request body.
func (i *Builder) WithBody(body any) *Builder {
	if i.Error != nil || body == nil {
		return i
	}

	var reqBody []byte
	var err error
	if reqBody, err = json.Marshal(body); err != nil {
		i.Error = err
		return i
	}

	i.HttpRequest.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	i.HttpRequest.Header.Set(string(ContentTypeKey), string(ApplicationJson))
	return i
}

// WithAuth sets the request `Authorization` header .
func (i *Builder) WithAuth(auth string) *Builder {
	if i.Error != nil {
		return i
	}
	i.HttpRequest.Header.Set(string(Authorization), auth)
	return i
}

// WithHeaders sets the request headers.
func (i *Builder) WithHeaders(headers map[string]string) *Builder {
	if i.Error != nil {
		return i
	}
	for k, v := range headers {
		i.HttpRequest.Header.Add(k, v)
	}
	return i
}

// Build returns the http request and error.
func (i *Builder) Build() (*http.Request, error) {
	return i.HttpRequest, i.Error
}
