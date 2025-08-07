package middleware

import (
	"github.com/go-kit/log"
	"github.com/minhho2511/elotusteam-test/utils"
	"net/http"
	"runtime/debug"
)

func MuxRecovery(logger log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					utils.ResponseWriter(w, http.StatusInternalServerError, utils.SetDefaultResponse(req.Context(), utils.Message{Code: 500}))
					_ = logger.Log("panic", err)
					debug.PrintStack()
				}
			}()
			h.ServeHTTP(w, req)
		})
	}
}

func Recovery(logger log.Logger) {
	if err := recover(); err != nil {
		_ = logger.Log("panic", err)
		debug.PrintStack()
	}
}
