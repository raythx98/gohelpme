package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwt"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"net/http"
	"strconv"
)

type AuthType string

const (
	Access  AuthType = "access"
	Refresh AuthType = "refresh"
	Basic   AuthType = "basic"
)

func Auth(authType AuthType) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			reqCtx := reqctx.GetValue(r.Context())
			if authType == Access {
				token, err := request.BearerExtractor{}.ExtractToken(r)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: err})
					return
				}

				jwtToken, err := jwt.GetValidAccessToken(token)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid access token")})
					return
				}

				subject, err := jwtToken.Claims.GetSubject()
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid subject")})
					return
				}

				parseInt, err := strconv.ParseInt(subject, 10, 64)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("failed to parse subject")})
					return
				}

				reqCtx.SetUserId(parseInt)
			}

			if authType == Refresh {
				token, err := request.BearerExtractor{}.ExtractToken(r)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: err})
					return
				}

				jwtToken, err := jwt.GetValidRefreshToken(token)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid refresh token")})
					return
				}

				subject, err := jwtToken.Claims.GetSubject()
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid subject")})
					return
				}

				parseInt, err := strconv.ParseInt(subject, 10, 64)
				if err != nil {
					reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("failed to parse subject")})
					return
				}

				reqCtx.SetUserId(parseInt)
			}

			if authType == Basic {
				// Pass
			}

			next.ServeHTTP(w, r)
		}
	}
}
