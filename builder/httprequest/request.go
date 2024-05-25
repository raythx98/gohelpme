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

func (i *Builder) New(ctx context.Context, method Method, url string) *Builder {
	i.HttpRequest, i.Error = http.NewRequestWithContext(ctx, string(method), url, nil)
	return i
}

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

func (i *Builder) WithAuth(auth string) *Builder {
	if i.Error != nil {
		return i
	}
	i.HttpRequest.Header.Set(string(Authorization), auth)
	return i
}

func (i *Builder) WithHeaders(headers map[string]string) *Builder {
	if i.Error != nil {
		return i
	}
	for k, v := range headers {
		i.HttpRequest.Header.Add(k, v)
	}
	return i
}

func (i *Builder) Build() (*http.Request, error) {
	return i.HttpRequest, i.Error
}
