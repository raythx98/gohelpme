package middleware

import (
	"context"
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

// ReqCtx is a middleware that adds a request context to the request.
func ReqCtx(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newCtx := context.WithValue(r.Context(), reqctx.Key, reqctx.New(r.Header.Get(string(httprequest.RequestId))))
		next.ServeHTTP(w, r.WithContext(newCtx))
	}
}
