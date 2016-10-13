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
	ctx := r.Context()
	s := ctx.Value("pluto").(Service)
	hcr := s.Health()
	if hcr.Status.String() != "SERVING" {
		reply.Json(w, r, http.StatusTooManyRequests, hcr)
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func newHealthServer() server.Server {
	// Define Router
	mux := router.NewMux()
	mux.GET("/_health", healthHandler)
	// Define server
	return server.NewServer(
		server.Name("health"),
		server.Addr(":9090"),
		server.Mux(mux))
}

func (s *service) startHealthHTTPServer() {
	s.health.SetServingStatus(s.cfg.Name, 1)
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
	s.health.SetServingStatus(s.cfg.Name, 2)
	// add go routine to WaitGroup
	s.wg.Add(1)
	go func(srv server.Server) {
		defer s.wg.Done()
		srv.Stop()
	}(s.healthHTTP)
}

func (s *service) healthOnServers() *healthpb.HealthCheckResponse {
	var hrc = &healthpb.HealthCheckResponse{Status: 1}
	for _, srv := range s.cfg.Servers {
		h := srv.Health()
		s.health.SetServingStatus(srv.Config().Name, h.Status)
		if h.Status.String() != "SERVING" {
			hrc.Status = h.Status
		}
	}
	return hrc
}

func (s *service) healthOnClients() *healthpb.HealthCheckResponse {
	var hrc = &healthpb.HealthCheckResponse{Status: 1}
	for _, clt := range s.cfg.Clients {
		h := clt.Health()
		s.health.SetServingStatus(clt.Config().Name, h.Status)
		if h.Status.String() != "SERVING" {
			hrc.Status = h.Status
		}
	}
	return hrc
}
