package middleware

import (
	"net/http"

	"github.com/raythx98/gohelpme/builder/httprequest"
)

func JsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(string(httprequest.ContentTypeKey), string(httprequest.ApplicationJson))
			next.ServeHTTP(w, r)
		},
	)
}
