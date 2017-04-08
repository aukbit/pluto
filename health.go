package pluto

import (
	"net/http"
	"strings"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var hcr = &healthpb.HealthCheckResponse{Status: 0}
	ctx := r.Context()
	m := ctx.Value(router.Key("module")).(string)
	n := ctx.Value(router.Key("name")).(string)
	s := ctx.Value(Key("pluto")).(*Service)

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
	case "db":
		db, err := s.Datastore()
		if err != nil {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		if n != db.Name() {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = db.Health()
	case "pluto":
		if n != s.Config().Name {
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
