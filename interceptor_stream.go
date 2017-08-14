package pluto

import (
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc"
)

// serviceContextStreamServerInterceptor Interceptor that adds service instance
// available in handlers context
func serviceContextStreamServerInterceptore(s *Service) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		// Note: service instance is always available in handlers context
		// under the general name > pluto
		ctx = s.WithContext(ctx)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
