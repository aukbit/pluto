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

// eidMiddleware sets eid in incoming metadata context
func eidMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			eidHeader := "X-Pluto-Eid"
			if _, ok := r.Header[eidHeader]; !ok {
				eid := common.RandID("", 16)
				r.Header.Add(eidHeader, eid)
			}
			if _, ok := w.Header()[eidHeader]; !ok {
				w.Header().Add(eidHeader, r.Header.Get(eidHeader))
			}
			ctx := r.Context()
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				md = metadata.New(map[string]string{})
			}
			md = md.Copy()
			md = metadata.Join(md, metadata.Pairs("eid", r.Header.Get(eidHeader)))
			ctx = metadata.NewIncomingContext(ctx, md)
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
			e := eidFromIncomingContext(ctx)
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

// eidFromIncomingContext returns eid from incoming context
func eidFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s := FromContext(ctx)
		l := s.Logger()
		l.Warn().Msg("metadata not available in incoming context")
		return ""
	}
	_, ok = md["eid"]
	if !ok {
		s := FromContext(ctx)
		l := s.Logger()
		l.Warn().Msg("eid not available in metadata")
		return ""
	}
	return md["eid"][0]
}

// eidFromOutgoingContext returns eid from outgoing context
func eidFromOutgoingContext(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		s := FromContext(ctx)
		l := s.Logger()
		l.Warn().Msg(fmt.Sprintf("%s metadata not available in outgoing context", s.Name()))
		return ""
	}
	_, ok = md["eid"]
	if !ok {
		s := FromContext(ctx)
		l := s.Logger()
		l.Warn().Msg(fmt.Sprintf("%s eid not available in metadata", s.Name()))
		return ""
	}
	return md["eid"][0]
}
