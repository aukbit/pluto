package client

import (
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LoggerUnaryInterceptor ...
func LoggerUnaryInterceptor(l zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromContext(ctx)
		if ok {
			e, ok := md["event"]
			if ok {
				l.Info("call", zap.String("event", e[0]), zap.String("method", method))
			}
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
