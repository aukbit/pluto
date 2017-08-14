package client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientStreamWithContext struct {
	cs grpc.ClientStream
	// ctx is the Context which we can assign it.
	ctx context.Context
}

func (w *ClientStreamWithContext) SetContext(ctx context.Context) {
	w.ctx = ctx
}

// Stream interface
func (w *ClientStreamWithContext) Context() context.Context      { return w.ctx }
func (w *ClientStreamWithContext) RecvMsg(msg interface{}) error { return w.cs.RecvMsg(msg) }
func (w *ClientStreamWithContext) SendMsg(msg interface{}) error { return w.cs.SendMsg(msg) }

// ClientStream interface
func (w *ClientStreamWithContext) Header() (metadata.MD, error) { return w.cs.Header() }
func (w *ClientStreamWithContext) Trailer() metadata.MD         { return w.cs.Trailer() }
func (w *ClientStreamWithContext) CloseSend() error             { return w.cs.CloseSend() }

// WrapClientStreamWithContext returns a ClientStream that has the ability to overwrite context.
func WrapClientStreamWithContext(stream grpc.ClientStream) *ClientStreamWithContext {
	exists, ok := stream.(*ClientStreamWithContext)
	if ok {
		return exists
	}
	return &ClientStreamWithContext{cs: stream, ctx: stream.Context()}
}

// WrapperStreamClient creates a single interceptor out of a chain of many interceptors
// Execution is done in right-to-left order
func WrapperStreamClient(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		h := wrapStream(streamer, interceptors...)
		return h(ctx, desc, cc, method, opts...)
	}
}

// wrap h with all specified interceptors
func wrapStream(ui grpc.Streamer, interceptors ...grpc.StreamClientInterceptor) grpc.Streamer {
	for _, i := range interceptors {
		h := func(current grpc.StreamClientInterceptor, next grpc.Streamer) grpc.Streamer {
			return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				return current(ctx, desc, cc, method, next, opts...)
			}
		}
		ui = h(i, ui)
	}
	return ui
}
