package pluto

import (
	"net/http"
	"strings"

	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func readyHealthHandler(w http.ResponseWriter, r *http.Request) {
	var hcr = &healthpb.HealthCheckResponse{Status: 0}
	ctx := r.Context()
	s := FromContext(ctx)
	// Test all servers
	for _, srv := range s.cfg.Servers {
		hcr = srv.Health()
		if hcr.Status.String() != healthpb.HealthCheckResponse_SERVING.String() {
			reply.Json(w, r, http.StatusServiceUnavailable, hcr)
			return
		}
	}
	// Test all clients
	for _, clt := range s.cfg.Clients {
		hcr = clt.Health()
		if hcr.Status.String() != healthpb.HealthCheckResponse_SERVING.String() {
			reply.Json(w, r, http.StatusServiceUnavailable, hcr)
			return
		}
	}
	// Test service
	hcr = s.Health()
	if hcr.Status.String() != healthpb.HealthCheckResponse_SERVING.String() {
		reply.Json(w, r, http.StatusServiceUnavailable, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)

}

func liveHealthHandler(w http.ResponseWriter, r *http.Request) {
	var hcr = &healthpb.HealthCheckResponse{Status: 0}
	ctx := r.Context()
	s := FromContext(ctx)
	if len(s.cfg.Servers) == 0 {
		reply.Json(w, r, http.StatusServiceUnavailable, hcr)
		return
	}
	// Test all servers
	for _, srv := range s.cfg.Servers {
		hcr = srv.Health()
		if hcr.Status.String() != healthpb.HealthCheckResponse_SERVING.String() {
			reply.Json(w, r, http.StatusServiceUnavailable, hcr)
			return
		}
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var hcr = &healthpb.HealthCheckResponse{Status: 0}
	ctx := r.Context()
	m := router.FromContext(ctx, "module")
	n := router.FromContext(ctx, "name")
	s := FromContext(ctx)

	switch m {
	case "server":
		name := strings.Replace(n, "_"+server.DefaultName, "", 1)
		srv, ok := s.Server(name)
		if !ok {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = srv.Health()
	case "client":
		name := strings.Replace(n, "_"+client.DefaultName, "", 1)
		clt, ok := s.Client(name)
		if !ok {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = clt.Health()
	case "pluto":
		if n != s.Name() {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = s.Health()
	}
	if hcr.Status.String() != healthpb.HealthCheckResponse_SERVING.String() {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}
