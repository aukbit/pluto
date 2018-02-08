package client

// loggerUnaryClientInterceptor ...
import (
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

func dialStreamClientInterceptor(clt *Client) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		e, ctx := eidFromContext(ctx)
		// sets new logger instance with eventID
		sublogger := clt.logger.With().Str("eid", e).Logger()
		sublogger.Info().Str("method", method).Msg(fmt.Sprintf("request %s", method))
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
