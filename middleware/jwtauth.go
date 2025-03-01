package middleware

import (
	"fmt"
	"net/http"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

// AuthType is the type of authentication.
// It can be either AccessToken or RefreshToken.
type AuthType string

const (
	AccessToken  AuthType = "AccessToken"
	RefreshToken AuthType = "RefreshToken"
)

// JwtAuth is a middleware that authenticates the request using JWT.
//
// It requires a jwthelper.IJwt implementation to authenticate the request.
func JwtAuth(jwtHelper jwthelper.IJwt, authType AuthType) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				if err != nil {
					reqctx.GetValue(r.Context()).SetError(errorhelper.NewAuthError(err))
				}
			}()

			switch authType {
			case AccessToken:
				if err = AuthenticateAccessToken(r, jwtHelper); err != nil {
					return
				}
			case RefreshToken:
				if err = AuthenticateRefreshToken(r, jwtHelper); err != nil {
					return
				}
			default:
				err = fmt.Errorf("invalid auth type")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func AuthenticateAccessToken(r *http.Request, jwtHelper jwthelper.IJwt) error {
	if err := jwtHelper.Authenticate(r, jwthelper.AccessToken); err != nil {
		return fmt.Errorf("%v, invalid access token", err)
	}

	return nil
}

func AuthenticateRefreshToken(r *http.Request, jwtHelper jwthelper.IJwt) error {
	if err := jwtHelper.Authenticate(r, jwthelper.AccessToken); err != nil {
		return fmt.Errorf("%v, invalid refresh token", err)
	}

	return nil
}
