package client

import (
	"github.com/uber-go/zap"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func (dc *defaultClient) dialGRPC() error {
	dc.logger.Info("dial")
	// establishes gRPC client connection
	// TODO use TLS
	// append logger
	dc.cfg.UnaryClientInterceptors = append(dc.cfg.UnaryClientInterceptors, loggerUnaryClientInterceptor(dc))
	// dial to establish connection
	conn, err := grpc.Dial(
		dc.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(WrapperUnaryClient(dc.cfg.UnaryClientInterceptors...)))

	if err != nil {
		dc.logger.Error("dial", zap.String("err", err.Error()))
		return err
	}
	// keep connection for later close
	dc.conn = conn
	// register health methods on connection
	dc.health = healthpb.NewHealthClient(conn)
	// register methods on connection
	dc.call = dc.cfg.GRPCRegister(conn)
	return nil
}
