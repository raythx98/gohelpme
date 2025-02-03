package middleware

import (
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"
)

// JsonResponse adds the application/json content type to the response header.
func JsonResponse(next Handler) Handler {
	return HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Add(string(httprequest.ContentTypeKey), string(httprequest.ApplicationJson))
			return next.ServeHTTP(w, r)
		},
	)
}
