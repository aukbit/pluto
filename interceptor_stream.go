package pluto

import (
	"github.com/aukbit/pluto/datastore"
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
		// ctx = context.WithValue(ctx, contextKey("pluto"), s)
		ctx = s.WithContext(ctx)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}

// serviceContextUnaryServerInterceptor Interceptor that adds service instance
// available in handlers context
func datastoreContextStreamServerInterceptor(s *Service) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		// get datastore from pluto service
		db, err := s.Datastore()
		if err != nil {
			return handler(srv, ss)
		}
		// requests new db session
		session, err := db.NewSession()
		if err != nil {
			return err
		}
		defer db.Close(session) // clean up
		ctx = datastore.WithContextSession(ctx, session)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
