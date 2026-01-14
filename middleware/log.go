package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

type responseBodyWriter struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

// Write writes the response body to the underlying ResponseWriter.
//
// It also captures the response body for logging.
// This method satisfies the http.ResponseWriter interface.
func (w *responseBodyWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}

// WriteHeader writes the response status code to the underlying ResponseWriter.
//
// It also captures the response status code for logging.
// This method satisfies the http.ResponseWriter interface.
func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LogOptions defines the options for the Log middleware.
type LogOptions struct {
	RequestRedact  []string
	ResponseRedact []string
}

// Log is a middleware that logs the request and response to the logger.
// It also logs some metadata about the request, such as the request ID, source IP, and execution time.
func Log(log logger.ILogger, opts ...LogOptions) func(next http.HandlerFunc) http.HandlerFunc {
	var opt LogOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			startAt := time.Now()

			// capture request body
			requestBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// capture response body
			respWriter := responseBodyWriter{w, make([]byte, 0), http.StatusOK}

			next.ServeHTTP(&respWriter, r)

			timeTaken := time.Since(startAt).String()

			logs := map[string]interface{}{
				"hostname":       getHostname(),
				"remote address": r.RemoteAddr,
				"request": map[string]interface{}{
					"started at": startAt.Truncate(time.Second),
					"endpoint":   fmt.Sprintf("%s %s://%s%s %s", r.Method, httphelper.GetScheme(r), r.Host, r.RequestURI, r.Proto),
					"headers":    r.Header,
					"body":       string(requestBody),
				},
				"response": map[string]interface{}{
					"completed at": time.Now().Truncate(time.Second),
					"status code":  respWriter.statusCode,
					"body":         string(respWriter.body),
				},
			}

			if len(opt.RequestRedact) > 0 {
				redact(logs["request"].(map[string]interface{}), opt.RequestRedact)
			}
			if len(opt.ResponseRedact) > 0 {
				redact(logs["response"].(map[string]interface{}), opt.ResponseRedact)
			}

			if reqctx.GetValue(r.Context()).Error != nil {
				log.Error(r.Context(), fmt.Sprintf("[in-http] %s %s%s %d in %s",
					r.Method, r.Host, r.RequestURI, respWriter.statusCode, timeTaken),
					logger.WithFields(logs),
				)
			} else {
				log.Info(r.Context(), fmt.Sprintf("[in-http] %s %s%s %d in %s",
					r.Method, r.Host, r.RequestURI, respWriter.statusCode, timeTaken),
					logger.WithFields(logs),
				)
			}
		}
	}
}

func redact(data map[string]interface{}, paths []string) {
	for _, path := range paths {
		parts := strings.Split(path, ".")
		redactPath(data, parts)
	}
}

func redactPath(data interface{}, parts []string) interface{} {
	if len(parts) == 0 {
		return data
	}

	key := parts[0]

	// If data is a string, it might be a JSON string
	if s, ok := data.(string); ok {
		var jsonData interface{}
		if err := json.Unmarshal([]byte(s), &jsonData); err == nil {
			redacted := redactPath(jsonData, parts)
			if b, err := json.Marshal(redacted); err == nil {
				return string(b)
			}
		}
		// If not JSON or error, we can't go deeper
		return data
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		// Might be a slice
		if l, ok := data.([]interface{}); ok {
			for i, v := range l {
				l[i] = redactPath(v, parts)
			}
			return l
		}
		return data
	}

	if len(parts) == 1 {
		if _, exists := m[key]; exists {
			m[key] = "*REDACTED*"
		}
		return m
	}

	if next, exists := m[key]; exists {
		m[key] = redactPath(next, parts[1:])
	}
	return m
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
