package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/raythx98/gohelpme/tool/validator"
	"io"
	"net/http"
)

// GetScheme returns the scheme of the request.
func GetScheme(req *http.Request) any {
	if req.TLS == nil {
		return "http"
	}
	return "https"
}

// CopyRequestBody copies the request body and returns it as a string.
//
// It also sets the request body back to its original state.
func CopyRequestBody(req *http.Request) string {
	requestBody := make([]byte, 0)
	if req.Body != nil {
		requestBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}
	return string(requestBody)
}

// CopyResponseBody copies the response body and returns it as a string.
//
// It also sets the response body back to its original state.
func CopyResponseBody(resp *http.Response) string {
	responseBody := make([]byte, 0)
	if resp.Body != nil {
		responseBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	}
	return string(responseBody)
}

// GetRequestBodyAndValidate reads the request body and validates it.
//
// It returns the request body and an error if any.
func GetRequestBodyAndValidate[T any](ctx context.Context, r *http.Request, v validator.IValidator) (T, error) {
	var body T
	requestByte, err := io.ReadAll(r.Body)
	if err != nil {
		return body, fmt.Errorf("failed to read request body: %w", err)
	}
	if err := json.Unmarshal(requestByte, &body); err != nil {
		return body, fmt.Errorf("failed to unmarshal request body: %w", err)
	}
	
	if err := v.StructCtx(ctx, body); err != nil {
		return body, fmt.Errorf("validation failed: %w", err)
	}
	return body, nil
}
