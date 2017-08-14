package ext

import (
	"github.com/aukbit/pluto/server"
	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	mgo "gopkg.in/mgo.v2"
)

// CassandraStreamServerInterceptor creates a Cassandra session and wraps it to context
func CassandraStreamServerInterceptor(name string, cfg gocql.ClusterConfig) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		s, err := gocql.NewSession(cfg)
		defer s.Close()
		if err != nil {
			zerolog.Ctx(ctx).Error().Msg(err.Error())
			return handler(srv, ss)
		}
		ctx = WithContextAny(ctx, name, s)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}

// MongoDBStreamServerInterceptor reates a MongoDB session and wraps it to context
func MongoDBStreamServerInterceptor(name string, cfg mgo.DialInfo) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		s, err := mgo.DialWithInfo(&cfg)
		defer s.Close()
		if err != nil {
			zerolog.Ctx(ctx).Error().Msg(err.Error())
			return handler(srv, ss)
		}
		ctx = WithContextAny(ctx, name, s)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}

// RedisStreamServerInterceptor wrap redis client to grpc
func RedisStreamServerInterceptor(name string, clt *redis.Client) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		ctx = WithContextAny(ctx, name, clt)
		// wrap context
		wrapped := server.WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
