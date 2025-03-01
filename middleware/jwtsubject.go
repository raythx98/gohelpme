package middleware

import (
	"fmt"
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

			if subjectString, err := jwtHelper.GetSubject(r); err == nil {
				fmt.Println("subjectString", subjectString)
				if subject, err := strconv.ParseInt(subjectString, 10, 64); err == nil {
					fmt.Println("subject", subject)
					reqctx.GetValue(r.Context()).SetUserId(subject)
				}
			}

			next.ServeHTTP(w, r)
		}
	}
}
