package server

import (
	"net/http"

	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/google/uuid"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// getOrCreateEventID uses grpc metadata context to set an event id
// the metadata context is then sent over the wire - gRPC calls
// and available to other services
func getOrCreateEventID(ctx context.Context) (string, context.Context) {
	// get
	md, ok := metadata.FromContext(ctx)
	if ok {
		e, ok := md["event"]
		if ok {
			return e[0], ctx
		}
	}
	// create
	e := uuid.New().String()
	ctx = metadata.NewContext(ctx, metadata.Pairs("event", e))
	return e, ctx
}

// middlewareServer Middleware to wrap all handlers with a server logger
func middlewareServer(srv *defaultServer) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// get or create unique event id for every request
			e, ctx := getOrCreateEventID(r.Context())
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

// middlewareStrictSecurityHeader Middleware to wrap all handlers with
// Strict-Transport-Security header
func middlewareStrictSecurityHeader() router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			h.ServeHTTP(w, r)
		}
	}
}
