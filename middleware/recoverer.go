package middleware

import (
	"fmt"
	"github.com/raythx98/gohelpme/tool/logger"
	"net/http"
)

// Recoverer recovers from panics and returns a 500 Internal Server Error.
func Recoverer(log logger.ILogger) func(next Handler) Handler {
	return func(next Handler) Handler {
		return HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) error {
				var err error
				defer func() {
					if p := recover(); p != nil {
						err = fmt.Errorf("panic: %v", p)
					}
				}()
				err = next.ServeHTTP(w, r)
				return err
			},
		)
	}
}
