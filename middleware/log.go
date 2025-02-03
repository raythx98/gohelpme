package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raythx98/gohelpme/builder/httprequest"
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

// Log is a middleware that logs the request and response to the logger.
// It also logs some metadata about the request, such as the request ID, source IP, and execution time.
func Log(log logger.ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				startAt := time.Now()
				ctx := context.WithValue(r.Context(), reqctx.Key, reqctx.New(r.Header.Get(string(httprequest.RequestId))))
				r = r.WithContext(ctx)

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
			},
		)
	}
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
