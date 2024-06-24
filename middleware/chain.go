package middleware

import (
	"net/http"
	"slices"
)

// Chain is a variadic function that takes a Handler and a list of functions that return a Handler.
// It returns a Handler that chains the functions together in the order they are provided.
// The first function in the list is the outermost function, which is applied first.
// And the last function in the list is the innermost function, which is applied last.
//
// Example:
//
//	defaultMiddlewares := []func(next http.Handler) http.Handler{
//	middleware.JsonResponse,
//		middleware.AddRequestId,
//		middleware.Log,
//  }
//	http.NewServeMux().Handle("/endpoint", middleware.Chain(method, defaultMiddlewares...))
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
