package server

import (
	"fmt"

	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (ds *defaultServer) healthGRPC() {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", ds.cfg.Port()), grpc.WithInsecure())
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.Name, 2)
	}
	defer conn.Close()
	c := healthpb.NewHealthClient(conn)
	hcr, err := c.Check(context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.Name})
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		ds.health.SetServingStatus(ds.cfg.Name, 2)
	}
	ds.health.SetServingStatus(ds.cfg.Name, hcr.Status)
}
