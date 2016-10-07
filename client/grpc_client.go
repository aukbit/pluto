package client

import (
	"github.com/uber-go/zap"

	"google.golang.org/grpc"
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
	// register methods on connection
	dc.call = dc.cfg.GRPCRegister(conn)
	return nil
}
