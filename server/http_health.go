package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (ds *defaultServer) healthHTTP() {
	r, err := http.Get(fmt.Sprintf(`http://localhost:%d/_health`, ds.cfg.Port()))
	if err != nil {
		// ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return
	}
	defer r.Body.Close()
	hcr := &healthpb.HealthCheckResponse{}
	if err := json.Unmarshal(b, hcr); err != nil {
		// ds.logger.Error("healthHttp", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return
	}
	ds.health.SetServingStatus(ds.cfg.ID, hcr.Status)
}
