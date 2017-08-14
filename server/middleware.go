package server

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/server/router"
	"github.com/rs/zerolog"
)

func serverMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := s.WithContext(r.Context())
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// eidMiddleware sets eid in outgoing metadata context
func eidMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			eid := common.RandID("", 16)
			ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("eid", common.RandID("", 16)))
			w.Header().Set("X-PLUTO-EID", eid)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// loggerMiddleware Middleware that adds logger instance
// available in handlers context and logs request
func loggerMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			e := eidFromOutgoingContext(ctx)
			// sets new logger instance with eid
			sublogger := s.logger.With().Str("eid", e).Logger()
			switch r.URL.Path {
			case "/_health":
				break
			default:
				h := zerolog.Dict()
				for k, v := range r.Header {
					h.Strs(k, v)
				}
				sublogger.Info().Str("method", r.Method).
					Str("url", r.URL.String()).
					Str("proto", r.Proto).
					Str("remote_addr", r.RemoteAddr).
					Dict("header", h).
					Msg(fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto))
			}
			// also nice to have a logger available in context
			ctx = sublogger.WithContext(ctx)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// strictSecurityHeaderMiddleware Middleware that adds
// Strict-Transport-Security header
func strictSecurityHeaderMiddleware() router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			h.ServeHTTP(w, r)
		}
	}
}

// --- Helper functions

func eidFromOutgoingContext(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		s := FromContext(ctx)
		s.Logger().Warn().Msg(fmt.Sprintf("%s metadata not available in outgoing context", s.Name()))
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