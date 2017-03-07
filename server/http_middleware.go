package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/server/router"
)

// loggerMiddleware Middleware that adds logger instance
// available in handlers context and logs request
func loggerMiddleware(srv *defaultServer) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// get or create unique event id for every request
			e, ctx := common.GetOrCreateEventID(r.Context())
			// create new log instance with eventID
			l := srv.logger.With(
				zap.String("event", e))
			l.Info("request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()))
			// also nice to have a logger available in context
			ctx = context.WithValue(ctx, "logger", l)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// strictSecurityHeaderMiddleware Middleware that adds
// Strict-Transport-Security header
func strictSecurityHeaderMiddleware() router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			h.ServeHTTP(w, r)
		}
	}
}
