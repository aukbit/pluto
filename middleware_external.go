package pluto

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/aukbit/pluto/server/router"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
)

// CassandraMiddleware creates a Cassandra session and add it to the context
func CassandraMiddleware(name string, cfg gocql.ClusterConfig) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s, err := gocql.NewSession(cfg)
			defer s.Close()
			if err != nil {
				zerolog.Ctx(r.Context()).Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := WithContextAny(r.Context(), name, s)
			// pass execution to the original handler
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// MongoDBMiddleware creates a MongoDB session and add it to the context
func MongoDBMiddleware(name string, cfg mgo.DialInfo) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			s, err := mgo.DialWithInfo(&cfg)
			defer s.Close()
			if err != nil {
				zerolog.Ctx(r.Context()).Error().Msg(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := WithContextAny(r.Context(), name, s)
			// pass execution to the original handler
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// RedisMiddleware adds Redis client to the context.
// Note: Caller is responsile to close the redis connection when its done.
// It is rare to Close a Client, as the Client is meant to be
// long-lived and shared between many goroutines.
func RedisMiddleware(name string, clt *redis.Client) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := WithContextAny(r.Context(), name, clt)
			// pass execution to the original handler
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
