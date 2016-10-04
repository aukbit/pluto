package server

import (
	"net/http"

	"bitbucket.org/aukbit/pluto/server/router"
)

// MiddlewareStrictSecurityHeader Middleware to wrap all handlers with
// Strict-Transport-Security header
func MiddlewareStrictSecurityHeader() router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			h.ServeHTTP(w, r)
		}
	}
}
