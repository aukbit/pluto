package datastore

import (
	"net/http"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"github.com/gocql/gocql"
)

// MiddlewareCassandra creates a db session and add it to the context
func MiddlewareDatastore(cluster *gocql.ClusterConfig) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get datastore from pluto service from context
			db := ctx.Value("pluto").(*pluto.Service).Config().Datastore
			// create db session
			session, err := db.Session()
			if err != nil {
				return reply.Json(w, r, http.StatusInternalServerError, err.Error())
			}
			defer session.Close() // clean up

			// save it in the router context
			context.Set(r, "dbsession", session)

			// pass execution to the original handler
			h.ServeHTTP(w, r)
		}
	}
}
