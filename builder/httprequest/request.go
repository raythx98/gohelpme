package httprequest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Builder interface {
	New(ctx context.Context, method string, url string) *Implementor
	WithBody(body any) *Implementor
	WithAuth(auth string) *Implementor
	WithHeaders(headers map[string]string) *Implementor
	Builder() (*http.Request, error)
}

type Implementor struct {
	HttpRequest *http.Request
	Error       error
}

func (i *Implementor) New(ctx context.Context, method Method, url string) *Implementor {
	i.HttpRequest, i.Error = http.NewRequestWithContext(ctx, string(method), url, nil)
	return i
}

func (i *Implementor) WithBody(body any) *Implementor {
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

func (i *Implementor) WithAuth(auth string) *Implementor {
	if i.Error != nil {
		return i
	}
	i.HttpRequest.Header.Set(string(Authorization), auth)
	return i
}

func (i *Implementor) WithHeaders(headers map[string]string) *Implementor {
	if i.Error != nil {
		return i
	}
	for k, v := range headers {
		i.HttpRequest.Header.Add(k, v)
	}
	return i
}

func (i *Implementor) Build() (*http.Request, error) {
	return i.HttpRequest, i.Error
}
