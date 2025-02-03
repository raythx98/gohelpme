package middleware

import "net/http"

// CORS adds the necessary headers to allow CORS requests.
func CORS(next Handler) Handler {
	return HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) error {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return nil
			}

			return next.ServeHTTP(w, r)
		},
	)
}
