package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/raythx98/gohelpme/util/reqctx"
	"github.com/raythx98/gohelpme/util/slogger"
)

func init() {
	slogger.Init()
}

type responseBodyWriter struct {
	http.ResponseWriter
	body []byte
}

func (w *responseBodyWriter) Write(body []byte) (int, error) {
	w.body = body
	return w.ResponseWriter.Write(body)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := slogger.AppendCtx(r.Context(), slog.String("request id", r.Header.Get(requestIdHeaderKey)))
			ctx = context.WithValue(ctx, reqctx.Key, reqctx.Value{RequestId: r.Header.Get(requestIdHeaderKey)})

			requestBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			bodyWriter := responseBodyWriter{w, make([]byte, 0)}

			next.ServeHTTP(&bodyWriter, r.WithContext(ctx))

			slog.LogAttrs(
				ctx, slog.LevelInfo, "request",
				slog.String("endpoint", fmt.Sprintf("%s %s%s %s", r.Method, r.Host, r.RequestURI, r.Proto)),
				slog.String("remote address", r.RemoteAddr),
				slog.Any("headers", r.Header),
				slog.String("body", string(requestBody)),
				slog.String("time taken", time.Since(start).String()),
				slog.String("response body", string(bodyWriter.body)),
			)
		},
	)
}
