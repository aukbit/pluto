package client

// loggerUnaryClientInterceptor ...
import (
	"github.com/aukbit/pluto/common"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func dialUnaryClientInterceptor(clt *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// get or create unique event id for every request
		e, ctx := common.GetOrCreateEventID(ctx)
		// create new log instance with eventID
		l := clt.logger.With(
			zap.String("event", e),
		)
		l.Info("dial",
			zap.String("method", method),
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
