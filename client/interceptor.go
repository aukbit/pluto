package client

// loggerUnaryClientInterceptor ...
import (
	"fmt"

	"github.com/aukbit/pluto/common"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func dialUnaryClientInterceptor(clt *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// get or create unique event id for every request
		e, ctx := common.GetOrCreateEventID(ctx)
		fmt.Println("client", ctx)
		// sets new logger instance with eventID
		sublogger := clt.logger.With().Str("event", e).Logger()
		sublogger.Info().Str("method", method).Msg(fmt.Sprintf("%s called %s", clt.Name(), method))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
