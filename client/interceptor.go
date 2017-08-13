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
		// set log instance with eventID
		clt.logger.With().Str("event", e)
		clt.logger.Info().Str("method", method).Msg(fmt.Sprintf("%s called %s", clt.Name(), method))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
