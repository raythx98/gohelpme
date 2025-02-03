package middleware

import (
	"net/http"
	"slices"
)

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// Chain is a variadic function that takes a Handler and a list of functions that return a Handler.
// It returns a Handler that chains the functions together in the order they are provided.
// The first function in the list is the outermost function, which is applied first.
// And the last function in the list is the innermost function, which is applied last.
//
// Example:
//
//		defaultMiddlewares := []func(next middleware.Handler) middleware.Handler{
//		middleware.JsonResponse,
//			middleware.AddRequestId,
//			middleware.Log,
//	 }
//		http.NewServeMux().Handle("/endpoint", middleware.Chain(method, defaultMiddlewares...))
func Chain(f HandlerFunc, m ...func(Handler) Handler) Handler {
	middlewares := slices.Clone(m)
	slices.Reverse(middlewares)

	var finalHandler Handler
	for _, candidate := range middlewares {
		if finalHandler == nil {
			finalHandler = candidate(f)
			continue
		}

		finalHandler = candidate(finalHandler)
	}
	return finalHandler

}
