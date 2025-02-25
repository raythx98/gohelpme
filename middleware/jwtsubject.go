package middleware

import (
	"net/http"
	"strconv"

	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

// JwtSubject is a middleware that extracts the subject from the JWT token and sets it in the request context.
//
// It requires a jwthelper.IJwt implementation to authenticate the request.
func JwtSubject(jwtHelper jwthelper.IJwt) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var err error
			reqCtx := reqctx.GetValue(r.Context())
			defer func() {
				reqCtx.SetError(err)
			}()

			subject, err := jwtHelper.GetSubject(r)
			if err != nil {
				return
			}

			parseInt, err := strconv.ParseInt(subject, 10, 64)
			if err != nil {
				return
			}

			reqCtx.SetUserId(parseInt)

			next.ServeHTTP(w, r)
		}
	}
}
