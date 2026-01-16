package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
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

func (w *responseBodyWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}

func (w *responseBodyWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type LogConfig struct {
	RedactedPaths []string
}

// Log is a middleware that logs the request and response to the logger.
func Log(log logger.ILogger, cfg LogConfig) func(next http.HandlerFunc) http.HandlerFunc {
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

			// redact sensitive information
			redactedLogs := redact(logs, cfg.RedactedPaths)

			if reqctx.GetValue(r.Context()).Error != nil {
				log.Error(r.Context(), fmt.Sprintf("[in-http] %s %s%s %d in %s",
					r.Method, r.Host, r.RequestURI, respWriter.statusCode, timeTaken),
					logger.WithFields(redactedLogs),
				)
			} else {
				log.Info(r.Context(), fmt.Sprintf("[in-http] %s %s%s %d in %s",
					r.Method, r.Host, r.RequestURI, respWriter.statusCode, timeTaken),
					logger.WithFields(redactedLogs),
				)
			}
		}
	}
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
