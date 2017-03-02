package pluto

import (
	"net/http"
	"strings"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var hcr = &healthpb.HealthCheckResponse{Status: 0}
	ctx := r.Context()
	m := ctx.Value("module").(string)
	n := ctx.Value("name").(string)
	s := ctx.Value("pluto").(Service)

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
		db, ok := s.Config().Datastore.(datastore.Datastore)
		if !ok {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		if n != db.Config().Name {
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
	if hcr.Status.String() != "SERVING" {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func newHealthServer(name string) server.Server {
	// Define Router
	mux := router.NewMux()
	mux.GET("/_health/:module/:name", healthHandler)
	// Define server
	return server.NewServer(
		server.Name(name+"_health"),
		server.Addr(":9090"),
		server.Mux(mux))
}

func (s *service) setHealthServer() {
	s.health.SetServingStatus(s.cfg.ID, 1)
	srv := newHealthServer(s.cfg.Name)
	s.cfg.Servers[srv.Config().Name] = srv
}
