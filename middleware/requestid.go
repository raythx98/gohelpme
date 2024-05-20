package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(requestIdHeaderKey) == "" {
				r.Header.Set(requestIdHeaderKey, uuid.NewString())
			}

			next.ServeHTTP(w, r)
		},
	)
}
