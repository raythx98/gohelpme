package middleware

import (
	"github.com/raythx98/gohelpme/tool/logger"
	"net/http"
	"runtime/debug"
)

// Recoverer recovers from panics and returns a 500 Internal Server Error.
func Recoverer(log logger.ILogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if p := recover(); p != nil {
						log.Error(r.Context(), "[panic]",
							logger.WithField("detail", p),
							logger.WithField("stack", string(debug.Stack())),
						)
						http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
					}
				}()
				next.ServeHTTP(w, r)
			},
		)
	}
}
