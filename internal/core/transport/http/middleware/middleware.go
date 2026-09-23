package core_http_middleware

import "net/http"

// Middleware — это тип функции, которая принимает http.Handler и возвращает http.Handler.

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
