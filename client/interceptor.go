package client

// loggerUnaryClientInterceptor ...
import (
	"fmt"

	"github.com/aukbit/pluto/common"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func dialUnaryClientInterceptor(clt *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// get or create unique event id for every request
		e, ctx := eidFromOutgoingContext(ctx)
		// sets new logger instance with eventID
		sublogger := clt.logger.With().Str("event", e).Logger()
		sublogger.Info().Str("method", method).Msg(fmt.Sprintf("%s called %s", clt.Name(), method))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// --- Helper functions

// eidFromOutgoingContext returns eid and context with eid in context metadata
func eidFromOutgoingContext(ctx context.Context) (string, context.Context) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	// NOTE: modification should be made to copies of the returned MD.
	md = md.Copy()
	// verify if eid exists, if not generate new eid
	_, ok = md["eid"]
	if !ok {
		md = metadata.Join(md, metadata.Pairs("eid", common.RandID("", 16)))
		ctx = metadata.NewOutgoingContext(ctx, md)
		return md["eid"][0], ctx
	}
	return md["eid"][0], ctx
}
