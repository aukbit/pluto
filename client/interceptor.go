package client

// loggerUnaryClientInterceptor ...
import (
	"errors"
	"fmt"

	"github.com/aukbit/pluto/common"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	errEidNotAvailableOnIncomingContext = errors.New("eid not available on incoming context")
	errEidNotAvailableOnOutgoingContext = errors.New("eid not available on outgoing context")
)

func dialUnaryClientInterceptor(clt *Client) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		e, ctx := eidFromContext(ctx)
		// sets new logger instance with eventID
		sublogger := clt.logger.With().Str("eid", e).Logger()
		sublogger.Info().
			Str("method", method).
			Str("data", fmt.Sprintf("%v", req)).
			Msg(fmt.Sprintf("request %s", method))
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// --- Helper functions

// eidFromOutgoingContext returns metadata if eid is available in outgoing context
func eidInOutgoingContext(ctx context.Context) (string, metadata.MD) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", nil
	}
	if _, ok := md["eid"]; !ok {
		return "", nil
	}
	return md["eid"][0], md
}

// eidInIncomingContext returns metadata if eid is available in incoming context
func eidInIncomingContext(ctx context.Context) (string, metadata.MD) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", nil
	}
	if _, ok := md["eid"]; !ok {
		return "", nil
	}
	return md["eid"][0], md
}

// eidFromContext returns eid and context with eid in context metadata
func eidFromContext(ctx context.Context) (string, context.Context) {
	if eid, md := eidInIncomingContext(ctx); md != nil {
		md = md.Copy()
		ctx = metadata.NewOutgoingContext(ctx, md)
		return eid, ctx
	}
	if eid, md := eidInOutgoingContext(ctx); md != nil {
		return eid, ctx
	}
	// if not in context create new outgoing context with new eid
	md := metadata.New(map[string]string{})
	md = metadata.Join(md, metadata.Pairs("eid", common.RandID("", 16)))
	ctx = metadata.NewOutgoingContext(ctx, md)
	return md["eid"][0], ctx
}
