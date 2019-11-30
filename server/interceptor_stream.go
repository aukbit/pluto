package server

import (
	"fmt"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
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
		// get information from peer
		p, _ := peer.FromContext(ctx)
		e := eidFromIncomingContext(ctx)
		// sets new logger instance with eventID
		sublogger := s.logger.With().
			Str("eid", e).
			Str("method", info.FullMethod).Logger()
		sublogger.Info().
			Dict("peer", zerolog.Dict().
				Str("addr", fmt.Sprintf("%v", p.Addr)).
				Str("auth", fmt.Sprintf("%v", p.AuthInfo))).
			Msgf("call %s received from %v", info.FullMethod, p.Addr)
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		// wrap context
		wrapped := WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}

// InterfaceStreamServerInterceptor wraps any type to grpc stream server
func InterfaceStreamServerInterceptor(name string, val interface{}) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		ctx = WithContextAny(ctx, name, val)
		// wrap context
		wrapped := WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
