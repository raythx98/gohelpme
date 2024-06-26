package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/raythx98/gohelpme/tool/logger"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/raythx98/gohelpme/builder/httprequest"
	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

var (
	log logger.ILogger
)

func init() {
	// Initialize a fallback logger
	log = logger.NewDefault()
}

// RegisterLogger registers a logger to be used by the middleware.
// This should be called once at the beginning of the program before any middleware is used.
//
// If no logger is registered, a default (fallback) logger is used, but it is not recommended to use it.
// Example:
//
//	middleware.RegisterLogger(projectLogger)
func RegisterLogger(l logger.ILogger) {
	log = l
}

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
func Log(next http.Handler) http.Handler {
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

			log.Info(ctx, fmt.Sprintf("[in-http] %s %s%s %d in %s",
				r.Method, r.Host, r.RequestURI, respWriter.statusCode, timeTaken),
				logger.WithField("hostname", getHostname()),
				logger.WithField("remote address", r.RemoteAddr),
				logger.WithField("request", map[string]interface{}{
					"started at": startAt.Truncate(time.Second),
					"endpoint":   fmt.Sprintf("%s %s://%s%s %s", r.Method, httphelper.GetScheme(r), r.Host, r.RequestURI, r.Proto),
					"headers":    r.Header,
					"body":       string(requestBody),
				}),
				logger.WithField("response", map[string]interface{}{
					"completed at": time.Now().Truncate(time.Second),
					"status code":  respWriter.statusCode,
					"body":         string(respWriter.body),
				}),
			)

		},
	)
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
