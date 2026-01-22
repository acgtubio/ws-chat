package middleware

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	currentHandler := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		currentHandler = middlewares[i](currentHandler)
	}

	return currentHandler
}
