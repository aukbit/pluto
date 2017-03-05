package server

import (
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (ds *defaultServer) healthGRPC() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", ds.cfg.Port()), grpc.WithInsecure())
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return
	}
	defer conn.Close()
	c := healthpb.NewHealthClient(conn)
	hcr, err := c.Check(context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.ID})
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return
	}
	ds.health.SetServingStatus(ds.cfg.ID, hcr.Status)
}
