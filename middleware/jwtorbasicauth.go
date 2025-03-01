package middleware

import (
	"fmt"
	"net/http"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/basicauth"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

// JwtOrBasicAuth is a middleware that authenticates the request using JWT, and falls back to basic auth if the JWT token is invalid.
func JwtOrBasicAuth(basicAuthHelper basicauth.IAuth, jwtHelper jwthelper.IJwt, authType AuthType) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				if err != nil {
					reqctx.GetValue(r.Context()).SetError(errorhelper.NewAuthError(err))
				}
			}()

			isJwtAuthenticated := false
			switch authType {
			case AccessToken:
				if err = AuthenticateAccessToken(r, jwtHelper); err == nil {
					isJwtAuthenticated = true
				}
			case RefreshToken:
				if err = AuthenticateRefreshToken(r, jwtHelper); err == nil {
					isJwtAuthenticated = true
				}
			default:
				isJwtAuthenticated = false
			}

			if !isJwtAuthenticated {
				if basicErr := basicAuthHelper.Authenticate(r); basicErr != nil {
					err = fmt.Errorf("invalid jwt: %v, invalid basic: %v", err, basicErr)
					return
				} else {
					err = nil
				}
			}

			next.ServeHTTP(w, r)
		}
	}
}
