package server

import (
	"fmt"
	"net/http"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/server/router"
	"github.com/rs/zerolog"
)

// loggerMiddleware Middleware that adds logger instance
// available in handlers context and logs request
func loggerMiddleware(s *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get or create unique event id for every request
			e, ctx := common.GetOrCreateEventID(r.Context())
			fmt.Println("http", ctx)
			// sets new logger instance with eventID
			sublogger := s.logger.With().Str("event", e).Logger()
			switch r.URL.Path {
			case "/_health":
				break
			default:
				sublogger.Info().Msg(fmt.Sprintf("%v %v %v %v", s.Name(), r.Method, r.URL, r.Proto))
				if e := s.logger.Debug(); e.Enabled() {
					h := zerolog.Dict()
					for k, v := range r.Header {
						h.Strs(k, v)
					}
					s.logger.Debug().Str("method", r.Method).
						Str("url", r.URL.String()).
						Str("proto", r.Proto).
						Str("remote_addr", r.RemoteAddr).
						Dict("header", h).
						Msg(fmt.Sprintf("%v %v %v %v", s.Name(), r.Method, r.URL, r.Proto))
				}
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
