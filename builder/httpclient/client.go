package httpclient

import (
	"net/http"
)

type Client interface {
	Do() (*http.Response, error)
}

type Implementor struct {
	httpClient *http.Client
}

func NewUserClient() *Implementor {
	httpClient := &http.Client{}
	return &Implementor{httpClient: httpClient}
}

func (i *Implementor) Do(req *http.Request) (*http.Response, error) {
	return i.httpClient.Do(req)
}
