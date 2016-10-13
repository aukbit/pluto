package server

import (
	"net/http"

	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server/router"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	n := ctx.Value("name").(string)
	h := ctx.Value("health").(*health.Server)
	hcr, err := h.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: n})
	if err != nil {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func HealthMiddleware(hs *health.Server) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "health", hs)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
