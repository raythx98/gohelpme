package httphelper

import (
	"bytes"
	"io"
	"net/http"
)

func GetScheme(req *http.Request) any {
	if req.TLS == nil {
		return "http"
	}
	return "https"
}

func CopyRequestBody(req *http.Request) string {
	requestBody := make([]byte, 0)
	if req.Body != nil {
		requestBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}
	return string(requestBody)
}

func CopyResponseBody(resp *http.Response) string {
	responseBody := make([]byte, 0)
	if resp.Body != nil {
		responseBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	}
	return string(responseBody)
}
