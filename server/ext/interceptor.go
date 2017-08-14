package ext

import (
	context "golang.org/x/net/context"
	mgo "gopkg.in/mgo.v2"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// CassandraUnaryServerInterceptor creates a Cassandra session and add it to the context
func CassandraUnaryServerInterceptor(name string, cfg gocql.ClusterConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		s, err := gocql.NewSession(cfg)
		defer s.Close()
		if err != nil {
			zerolog.Ctx(ctx).Error().Msg(err.Error())
			return handler(ctx, req)
		}
		ctx = WithContextAny(ctx, name, s)
		return handler(ctx, req)
	}
}

// MongoDBUnaryServerInterceptor creates a MongoDB session and add it to the context
func MongoDBUnaryServerInterceptor(name string, cfg mgo.DialInfo) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		s, err := mgo.DialWithInfo(&cfg)
		defer s.Close()
		if err != nil {
			zerolog.Ctx(ctx).Error().Msg(err.Error())
			return handler(ctx, req)
		}
		ctx = WithContextAny(ctx, name, s)
		return handler(ctx, req)
	}
}

// RedisUnaryServerInterceptor wrap redis client to grpc
// Note: Caller is responsile to close the redis connection when its done.
// It is rare to Close a Client, as the Client is meant to be
// long-lived and shared between many goroutines.
func RedisUnaryServerInterceptor(name string, clt *redis.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = WithContextAny(ctx, name, clt)
		return handler(ctx, req)
	}
}
