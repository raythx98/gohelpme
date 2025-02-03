package middleware

import (
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"

	"github.com/google/uuid"
)

// AddRequestId adds a request ID to the request context if it doesn't already exist.
//
// This is useful for logging and tracing.
// It uses the request ID from the API Gateway event if it exists, otherwise it generates a new one
// and adds it to the request context.
func AddRequestId(next Handler) Handler {
	return HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			if r.Header.Get(string(httprequest.RequestId)) == "" {
				r.Header.Set(string(httprequest.RequestId), uuid.NewString())
			}

			return next.ServeHTTP(w, r)
		},
	)
}
