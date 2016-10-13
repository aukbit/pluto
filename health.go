package pluto

import (
	"net/http"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
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
		srv, ok := s.Server(n)
		if !ok {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = srv.Health()
	case "client":
		clt, ok := s.Client(n)
		if !ok {
			reply.Json(w, r, http.StatusNotFound, hcr)
			return
		}
		hcr = clt.Health()
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

func newHealthServer() server.Server {
	// Define Router
	mux := router.NewMux()
	mux.GET("/_health/:module/:name", healthHandler)
	// Define server
	return server.NewServer(
		server.Name("health"),
		server.Addr(":9090"),
		server.Mux(mux))
}

func (s *service) startHealthHTTPServer() {
	s.health.SetServingStatus(s.cfg.ID, 1)
	// add go routine to WaitGroup
	s.wg.Add(1)
	go func(srv server.Server) {
		defer s.wg.Done()
		if err := srv.Run(
			server.ParentID(s.cfg.ID),
			server.Middlewares(serviceContextMiddleware(s))); err != nil {
			s.logger.Error("Run()", zap.String("err", err.Error()))
		}
	}(s.healthHTTP)
}

func (s *service) stopHealthHTTPServer() {
	s.health.SetServingStatus(s.cfg.ID, 2)
	// add go routine to WaitGroup
	s.wg.Add(1)
	go func(srv server.Server) {
		defer s.wg.Done()
		srv.Stop()
	}(s.healthHTTP)
}
