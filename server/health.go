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
	s := FromContext(ctx)
	hcr, err := s.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: s.cfg.ID})
	if err != nil {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func serverMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := s.WithContext(r.Context())
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
