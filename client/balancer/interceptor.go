package balancer

// loggerUnaryClientInterceptor ...
import (
	"github.com/aukbit/pluto/common"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func loggerUnaryClientInterceptor(conn *connector) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// get or create unique event id for every request
		_, ctx = common.GetOrCreateEventID(ctx)
		// create new log instance with eventID
		// l := conn.logger.With(
		// 	zap.String("event", e))
		// l.Info("call",
		// 	zap.String("method", method))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
