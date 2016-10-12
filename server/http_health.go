package server

import (
	"log"
	"net/http"

	"bitbucket.org/aukbit/pluto/reply"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	hcr := &healthpb.HealthCheckResponse{Status: 1}
	log.Printf("healthHandler %v", hcr)
	reply.Json(w, r, http.StatusOK, hcr)
}
