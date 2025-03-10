package middleware

import (
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"
)

// JsonResponse adds the application/json content type to the response header.
func JsonResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(string(httprequest.ContentTypeKey), string(httprequest.ApplicationJson))
		next.ServeHTTP(w, r)
	}
}
