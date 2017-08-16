package server

import (
	"fmt"

	"github.com/rs/zerolog"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
			Msg(fmt.Sprintf("request %s from %v", info.FullMethod, p.Addr))

		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		h, err := handler(ctx, req)
		sublogger.Info().
			Str("data", fmt.Sprintf("%v", h)).
			Msg(fmt.Sprintf("response %s to %v", info.FullMethod, p.Addr))
		return h, err
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
