package pluto

import (
	"context"
	"net/http"

	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/server/router"
)

// serviceContextMiddleware Middleware that adds service instance
// available in handlers context
func serviceContextMiddleware(s *Service) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get context
			ctx := r.Context()
			// Note: service instance is always available in handlers context
			// under the general name > pluto
			ctx = context.WithValue(ctx, contextKey("pluto"), s)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// datastoreContextMiddleware creates a db session and add it to the context
func datastoreContextMiddleware(s *Service) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get datastore from pluto service
			db, err := s.Datastore()
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			// requests new db session
			session, err := db.NewSession()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer db.Close(session) // clean up
			// get context
			ctx := r.Context()
			ctx = datastore.WithContext(ctx, session)
			// // save it in the router context
			// ctx = context.WithValue(ctx, contextKey("session"), session)
			// pass execution to the original handler
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
