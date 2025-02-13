package middleware

import "net/http"

// Example is a simple example of how middleware can be used to wrap the request handling logic.
func Example(next http.HandlerFunc) http.HandlerFunc {
	// We wrap our anonymous function, and cast it to a httpbuilder.HandlerFunc
	// Because our function signature matches ServeHTTP(w, r), this allows
	// our function (type) to implicitly satisfy the httpbuilder.Handler interface.
	return func(w http.ResponseWriter, r *http.Request) {
		// Logic before - reading request values, putting things into the
		// request context, performing authentication

		// Important that we call the 'next' handler in the chain. If we don't,
		// then request handling will stop here.
		next.ServeHTTP(w, r)
		// Logic after - useful for logging, metrics, etc.
		//
		// It's important that we don't use the ResponseWriter after we've called the
		// next handler: we may cause conflicts when trying to write the response
	}
}
