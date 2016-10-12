package server

import (
	"fmt"

	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (ds *defaultServer) healthGRPC() *healthpb.HealthCheckResponse {
	var hcr = &healthpb.HealthCheckResponse{Status: 2}
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", ds.cfg.Port()), grpc.WithInsecure())
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		return hcr
	}
	defer conn.Close()

	c := healthpb.NewHealthClient(conn)
	hcr, err = c.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		ds.logger.Error("healthGRPC", zap.String("err", err.Error()))
		return hcr
	}
	return hcr
}
