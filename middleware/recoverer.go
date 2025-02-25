package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/raythx98/gohelpme/tool/reqctx"
)

// Recoverer recovers from panics and returns a 500 Internal Server Error.
func Recoverer() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					reqCtx := reqctx.GetValue(r.Context())
					reqCtx.SetError(fmt.Errorf("[panic] %v", p))
					reqCtx.SetErrorStack(debug.Stack())
				}
			}()
			next.ServeHTTP(w, r)
		}
	}
}
