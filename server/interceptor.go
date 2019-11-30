package server

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func serverUnaryServerInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = s.WithContext(ctx)
		return handler(ctx, req)
	}
}

func loggerUnaryServerInterceptor(s *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		// get information from peer
		p, _ := peer.FromContext(ctx)
		e := eidFromIncomingContext(ctx)
		// sets new logger instance with eid
		sublogger := s.logger.With().
			Str("eid", e).
			Str("method", info.FullMethod).Logger()
		sublogger.Info().
			Str("data", fmt.Sprintf("%v", req)).
			Dict("peer", zerolog.Dict().
				Str("addr", fmt.Sprintf("%v", p.Addr)).
				Str("auth", fmt.Sprintf("%v", p.AuthInfo))).
			Msgf("call %s received from %v", info.FullMethod, p.Addr)

		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		h, err := handler(ctx, req)
		end := time.Now()
		sublogger.Info().
			Str("data", fmt.Sprintf("%v", h)).
			Msgf("response %s sent to %v - duration: %v", info.FullMethod, p.Addr, end.Sub(start))
		return h, err
	}
}

// InterfaceUnaryServerInterceptor wraps any type to grpc unary server
func InterfaceUnaryServerInterceptor(name string, val interface{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = WithContextAny(ctx, name, val)
		return handler(ctx, req)
	}
}
