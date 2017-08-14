package server

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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
		sublogger.Info().Str("method", info.FullMethod).Msg(fmt.Sprintf("request %s", info.FullMethod))
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
		s.Logger().Warn().Msg("metadata not available in incoming context")
		return ""
	}
	_, ok = md["eid"]
	if !ok {
		s := FromContext(ctx)
		s.Logger().Warn().Msg("eid not available in metadata")
		return ""
	}
	return md["eid"][0]
}
