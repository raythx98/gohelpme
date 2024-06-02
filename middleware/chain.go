package middleware

import (
	"net/http"
	"slices"
)

func Chain(f http.HandlerFunc, m ...func(http.Handler) http.Handler) http.Handler {
	middlewares := slices.Clone(m)
	slices.Reverse(middlewares)

	var finalHandler http.Handler
	for _, candidate := range middlewares {
		if finalHandler == nil {
			finalHandler = candidate(f)
			continue
		}

		finalHandler = candidate(finalHandler)
	}
	return finalHandler
}
