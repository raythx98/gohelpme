package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/raythx98/gohelpme/builder/httprequest"
	"github.com/raythx98/gohelpme/util/httphelper"
	"github.com/raythx98/gohelpme/util/reqctx"
	"github.com/raythx98/gohelpme/util/slogger"
)

func init() {
	slogger.Init()
}

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

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			startAt := time.Now()
			ctx := context.WithValue(r.Context(), reqctx.Key, reqctx.New(r.Header.Get(string(httprequest.RequestId)), nil))
			r = r.WithContext(ctx)

			// capture request body
			requestBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// capture response body
			respWriter := responseBodyWriter{w, make([]byte, 0), http.StatusOK}

			next.ServeHTTP(&respWriter, r)

			slog.LogAttrs(
				ctx, slog.LevelInfo, "incoming http execution",
				slog.String("hostname", getHostname()),
				slog.String("remote address", r.RemoteAddr),
				slog.String("time taken", time.Since(startAt).String()),
				slog.Group("request",
					slog.Time("started at", startAt.Truncate(time.Second)),
					slog.String("endpoint", fmt.Sprintf("%s %s://%s%s %s",
						r.Method, httphelper.GetScheme(r), r.Host, r.RequestURI, r.Proto)),
					slog.Any("headers", r.Header),
					slog.String("body", string(requestBody)),
				),
				slog.Group("response",
					slog.Time("completed at", time.Now().Truncate(time.Second)),
					slog.Int("status code", respWriter.statusCode),
					slog.String("body", string(respWriter.body)),
				),
			)
		},
	)
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}
