package middleware

import (
	"github.com/go-kit/log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Adapt(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func TraceIdentifier() Middleware {
	return TraceIdentifierMiddleware
}

func AppMiddleware(handler http.Handler, lg log.Logger) http.Handler {
	return Adapt(
		handler,
		MuxRecovery(lg),
		TraceIdentifier(),
	)
}
