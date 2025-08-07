package middleware

import (
	"context"
	"github.com/matoous/go-nanoid/v2"
	"net/http"
)

const (
	TraceIDContextKey       = "KD-Trace-ID"
	TraceIDRequestHeaderKey = "X-Correlation-ID"
)

func TraceIdentifierMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceId := r.Header.Get(TraceIDRequestHeaderKey)
		if traceId == "" {
			traceId, _ = gonanoid.New()
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, TraceIDContextKey, traceId) // nolint
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
