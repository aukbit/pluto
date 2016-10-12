package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto/reply"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	hcr := &healthpb.HealthCheckResponse{Status: 1}
	reply.Json(w, r, http.StatusOK, hcr)
}

func (ds *defaultServer) healthHTTP() *healthpb.HealthCheckResponse {
	var hcr = &healthpb.HealthCheckResponse{Status: 2}
	r, err := http.Get(fmt.Sprintf(`http://localhost:%d/_health`, ds.cfg.Port()))
	if err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		return hcr
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
	}
	defer r.Body.Close()
	if err := json.Unmarshal(b, hcr); err != nil {
		ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		return hcr
	}
	return hcr
}
