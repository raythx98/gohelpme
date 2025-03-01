package middleware

import (
	"github.com/raythx98/gohelpme/tool/basicauth"
	"net/http"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

// BasicAuth is a middleware that authenticates the request using basic auth.
//
// It requires a basicauth.IAuth implementation to authenticate the request.
func BasicAuth(basicAuthHelper basicauth.IAuth) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				if err != nil {
					reqctx.GetValue(r.Context()).SetError(errorhelper.NewAuthError(err))
				}
			}()

			if err = basicAuthHelper.Authenticate(r); err != nil {
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
