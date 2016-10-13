package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server/router"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h := ctx.Value("health").(*health.Server)
	hcr, err := h.Check(
		context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		// TODO
		return
	}
	reply.Json(w, r, http.StatusOK, hcr)
}

func (ds *defaultServer) healthHTTP() {
	r, err := http.Get(fmt.Sprintf(`http://localhost:%d/_health`, ds.cfg.Port()))
	if err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.Name, 2)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.Name, 2)
		return
	}
	defer r.Body.Close()
	hcr := &healthpb.HealthCheckResponse{}
	if err := json.Unmarshal(b, hcr); err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.Name, 2)
		return
	}
	ds.health.SetServingStatus(ds.cfg.Name, hcr.Status)
}

func healthMiddleware(hs *health.Server) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "health", hs)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
