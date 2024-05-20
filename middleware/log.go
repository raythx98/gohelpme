package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/raythx98/gohelpme/util/reqctx"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Info("request", "Scheme", r.Proto)
			slog.Info("request", "Method", r.Method)
			slog.Info("request", "Host", r.Host)
			slog.Info("request", "RemoteAddr", r.RemoteAddr)
			slog.Info("request", "RequestUri", r.RequestURI)
			slog.Info("request", "URL", r.URL)
			slog.Info("request", "Headers", r.Header)

			ctx := context.WithValue(r.Context(), reqctx.Key, reqctx.Value{RequestId: r.Header.Get(requestIdHeaderKey)})
			slog.InfoContext(ctx, "request", "body", r.Body)

			next.ServeHTTP(w, r.WithContext(ctx))

			slog.InfoContext(ctx, "response", "body", w)
		},
	)
}
