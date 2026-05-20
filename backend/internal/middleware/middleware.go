package middleware

import "net/http"

type Middleware interface {
	Middleware(
		func(http.ResponseWriter, *http.Request),
	) func(http.ResponseWriter, *http.Request)
}

func GenerateHandler(
	core func(http.ResponseWriter, *http.Request),
	middlewares []Middleware,
) func(http.ResponseWriter, *http.Request) {
	final_handler := core

	for _, middleware := range middlewares {
		final_handler = middleware.Middleware(final_handler)
	}

	return final_handler
}
