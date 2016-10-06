package client

import (
	"bitbucket.org/aukbit/pluto/common"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// WrapperUnaryClient creates a single interceptor out of a chain of many interceptors
// Execution is done in right-to-left order
func WrapperUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		h := wrap(invoker, interceptors...)
		return h(ctx, method, req, reply, cc, opts...)
	}
}

// wrap h with all specified interceptors
func wrap(ui grpc.UnaryInvoker, interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryInvoker {
	for _, i := range interceptors {
		h := func(current grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return current(ctx, method, req, reply, cc, next, opts...)
			}
		}
		ui = h(i, ui)
	}
	return ui
}

// loggerUnaryClientInterceptor ...
func loggerUnaryClientInterceptor(clt *gRPCClient) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// get or create unique event id for every request
		e, ctx := common.GetOrCreateEventID(ctx)
		// create new log instance with eventID
		l := clt.logger.With(
			zap.String("event", e))
		l.Info("call",
			zap.String("method", method))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
