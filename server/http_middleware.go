package server

import (
	"net/http"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/server/router"
)

// loggerMiddleware Middleware that adds logger instance
// available in handlers context and logs request
func loggerMiddleware(srv *Server) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get or create unique event id for every request
			_, ctx := common.GetOrCreateEventID(r.Context())
			// create new log instance with eventID
			// l := srv.logger.With(
			// 	zap.String("event", e))
			switch r.URL.Path {
			case "/_health":
				break
			default:
				// l.Info(fmt.Sprintf("%v %v %v %v", srv.Name(), r.Method, r.URL, r.Proto),
				// 	zap.String("method", r.Method),
				// 	zap.String("url", r.URL.String()),
				// 	zap.String("proto", r.Proto),
				// 	zap.String("remote_addr", r.RemoteAddr),
				// 	zap.Any("header", r.Header),
				// )
			}
			// also nice to have a logger available in context
			// ctx = context.WithValue(ctx, Key("logger"), l)
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
