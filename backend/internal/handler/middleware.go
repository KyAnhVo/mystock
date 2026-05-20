package handler

import "net/http"

type Middleware interface {
	Middleware(
		func(http.ResponseWriter, *http.Request),
	) func(http.ResponseWriter, *http.Request)
}

// Generate a handler from a core http handler and a list
// of middlewares.
//
// The order of applying is
// `middlewares[0](middlewares[1](...middlewares[n-1](core)...))`
func GenerateHandler(
	core func(http.ResponseWriter, *http.Request),
	middlewares []Middleware,
) func(http.ResponseWriter, *http.Request) {
	final_handler := core

	for i := len(middlewares) - 1; i >= 0; i -= 1 {
		middleware := middlewares[i]
		final_handler = middleware.Middleware(final_handler)
	}

	return final_handler
}
