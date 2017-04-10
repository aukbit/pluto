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
		ctx = context.WithValue(ctx, Key("pluto"), s)
		return handler(ctx, req)
	}
}

// datastoreContextUnaryServerInterceptor Interceptor that adds service instance
// available in handlers context
func datastoreContextUnaryServerInterceptor(s *Service) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// get datastore from pluto service
		db, err := s.Datastore()
		if err != nil {
			return handler(ctx, req)
		}
		// requests new db session
		session, err := db.NewSession()
		if err != nil {
			return nil, err
		}
		defer db.Close(session) // clean up
		// save it in the router context
		ctx = context.WithValue(ctx, Key("session"), session)
		return handler(ctx, req)
	}
}
