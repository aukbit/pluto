package pluto

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// serviceContextUnaryServerInterceptor Interceptor that adds service instance
// available in handlers context
func serviceContextUnaryServerInterceptor(s *Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Note: service instance is always available in handlers context
		// under the general name > pluto
		ctx = context.WithValue(ctx, "pluto", s)
		return handler(ctx, req)
	}
}
