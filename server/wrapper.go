package server

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// WrapperUnaryServer creates a single interceptor out of a chain of many interceptors
// Execution is done in right-to-left order
func WrapperUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h := wrap(handler, info, interceptors...)
		return h(ctx, req)
	}
}

// wrap h with all specified interceptors
func wrap(uh grpc.UnaryHandler, info *grpc.UnaryServerInfo, interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryHandler {
	for _, i := range interceptors {
		h := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return current(ctx, req, info, next)
			}
		}
		uh = h(i, uh)
	}
	return uh
}
