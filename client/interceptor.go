package client

// loggerUnaryClientInterceptor ...
import (
	"errors"
	"fmt"
	"time"

	"github.com/aukbit/pluto/v6/common"

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
		start := time.Now()
		e := eidFromIncomingContext(ctx)
		ctx = eidToOutgoingContext(ctx, e)
		// sets new logger instance with eventID
		sublogger := clt.logger.With().Str("eid", e).Logger()
		sublogger.Info().
			Str("method", method).
			Str("data", fmt.Sprintf("%v", req)).
			Msgf("call %s", method)
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		err := invoker(ctx, method, req, reply, cc, opts...)
		end := time.Now()
		sublogger.Info().Msgf("response received %s - duration: %v", method, end.Sub(start))
		return err
	}
}

// --- Helper functions

// eidInIncomingContext returns eid from metadata in incoming context
func eidFromIncomingContext(ctx context.Context) string {
	// get eid from incoming context or generate new one
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return common.RandID("", 16)
	}
	if _, ok := md["eid"]; !ok {
		return common.RandID("", 16)
	}
	return md["eid"][0]
}

// eidToOutgoingContext returns context with eid in metadata in outgoing context
func eidToOutgoingContext(ctx context.Context, eid string) context.Context {
	// add eid to outgoing context
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md = md.Copy()
	md = metadata.Join(md, metadata.Pairs("eid", eid))
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
