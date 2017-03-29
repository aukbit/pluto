package datastore

import (
	"context"

	"google.golang.org/grpc"
)

// datastoreContextUnaryServerInterceptor Interceptor that adds service instance
// available in handlers context
func DatastoreContextUnaryServerInterceptor(s *Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Note: service instance is always available in handlers context
		// under the general name > pluto
		ctx = context.WithValue(ctx, "pluto", s)
		return handler(ctx, req)
	}
}
