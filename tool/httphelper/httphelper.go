package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
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

func GetRequestBodyAndValidate[T any](_ context.Context, r *http.Request, v *validator.Validate) (T, error) {
	var body T
	requestByte, err := io.ReadAll(r.Body)
	if err != nil {
		return body, fmt.Errorf("failed to read request body: %w", err)
	}
	if err := json.Unmarshal(requestByte, &body); err != nil {
		return body, fmt.Errorf("failed to unmarshal request body: %w", err)
	}
	if err := v.Struct(body); err != nil {
		return body, fmt.Errorf("validation failed: %w", err)
	}

	return body, nil
}
