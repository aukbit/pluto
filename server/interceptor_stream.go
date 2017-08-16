package server

import (
	"fmt"

	"google.golang.org/grpc"
)

func serverStreamServerInterceptor(s *Server) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		ctx = s.WithContext(ctx)
		// wrap context
		wrapped := WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}

func loggerStreamServerInterceptor(s *Server) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		e := eidFromIncomingContext(ctx)
		// sets new logger instance with eventID
		sublogger := s.logger.With().Str("eid", e).Logger()
		sublogger.Info().Str("method", info.FullMethod).Msg(fmt.Sprintf("request %s", info.FullMethod))
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		// wrap context
		wrapped := WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
