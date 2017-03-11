package server

import (
	"net/http"

	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"golang.org/x/net/context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ds := ctx.Value("server").(*Server)
	hcr, err := ds.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.ID})
	if err != nil {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func serverMiddleware(srv *Server) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "server", srv)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
