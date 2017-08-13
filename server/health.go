package server

import (
	"net/http"

	"github.com/aukbit/pluto/reply"
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
