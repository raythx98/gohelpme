package middleware

import "net/http"

func JsonResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeaderKey, "application/json")
			next.ServeHTTP(w, r)
		},
	)
}
