package middleware

import (
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"

	"github.com/google/uuid"
)

func AddRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(string(httprequest.RequestId)) == "" {
				r.Header.Set(string(httprequest.RequestId), uuid.NewString())
			}

			next.ServeHTTP(w, r)
		},
	)
}
