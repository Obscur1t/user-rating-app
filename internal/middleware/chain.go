package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(mux http.Handler, middlewares ...Middleware) http.Handler {
	handler := mux
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
