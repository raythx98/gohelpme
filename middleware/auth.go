package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwt"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"net/http"
)

type AuthType string

const (
	Access  AuthType = "access"
	Refresh AuthType = "refresh"
	Basic   AuthType = "basic"
)

func Auth(authType AuthType) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				reqCtx := reqctx.GetValue(r.Context())
				if authType == Access {
					token, err := request.BearerExtractor{}.ExtractToken(r)
					if err != nil {
						reqCtx.SetError(&errorhelper.AuthError{Err: err})
						return
					}

					if jwt.IsValidAccessToken(token) != nil {
						reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid access token")})
						return
					}
				}

				if authType == Refresh {
					token, err := request.BearerExtractor{}.ExtractToken(r)
					if err != nil {
						reqCtx.SetError(&errorhelper.AuthError{Err: err})
						return
					}

					if jwt.IsValidRefreshToken(token) != nil {
						reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid refresh token")})
						return
					}
				}

				if authType == Basic {
					// Pass
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
