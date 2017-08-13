package server

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func serverUnaryServerInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = s.WithContext(ctx)
		return handler(ctx, req)
	}
}

func loggerUnaryServerInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		e := eidFromIncomingContext(ctx)
		// sets new logger instance with eid
		sublogger := s.logger.With().Str("eid", e).Logger()
		sublogger.Info().Str("method", info.FullMethod).
			Msg(fmt.Sprintf("%s request %s", s.Name(), info.FullMethod))
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		return handler(ctx, req)
	}
}

// --- Helper functions

// eidFromIncomingContext returns eid from incoming context
func eidFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s := FromContext(ctx)
		s.Logger().Warn().Msg(fmt.Sprintf("%s metadata not available in incoming context", s.Name()))
		return ""
	}
	_, ok = md["eid"]
	if !ok {
		s := FromContext(ctx)
		s.Logger().Warn().Msg(fmt.Sprintf("%s eid not available in metadata", s.Name()))
		return ""
	}
	return md["eid"][0]
}
