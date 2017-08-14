package pluto

import (
	"net/http"

	"github.com/aukbit/pluto/server/router"
)

// serviceContextMiddleware Middleware that adds service instance
// available in handlers context
func serviceContextMiddleware(s *Service) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Note: service instance is always available in handlers context
			// under the general name > pluto
			ctx := s.WithContext(r.Context())
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
