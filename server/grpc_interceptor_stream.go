package server

import (
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// https://groups.google.com/forum/#!topic/grpc-io/Q88GQFTPF1o
type ServerStreamWithContext struct {
	ss grpc.ServerStream
	// ctx is the Context which we can assign it.
	ctx context.Context
}

func (w *ServerStreamWithContext) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *ServerStreamWithContext) Context() context.Context {
	return w.ctx
}
func (w *ServerStreamWithContext) RecvMsg(msg interface{}) error   { return w.ss.RecvMsg(msg) }
func (w *ServerStreamWithContext) SendMsg(msg interface{}) error   { return w.ss.SendMsg(msg) }
func (w *ServerStreamWithContext) SendHeader(md metadata.MD) error { return w.ss.SendHeader(md) }
func (w *ServerStreamWithContext) SetHeader(md metadata.MD) error  { return w.ss.SetHeader(md) }
func (w *ServerStreamWithContext) SetTrailer(md metadata.MD)       { w.ss.SetTrailer(md) }

// WrapServerStreamWrapper returns a ServerStream that has the ability to overwrite context.
func WrapServerStreamWithContext(stream grpc.ServerStream) *ServerStreamWithContext {
	exists, ok := stream.(*ServerStreamWithContext)
	if ok {
		return exists
	}
	return &ServerStreamWithContext{ss: stream, ctx: stream.Context()}
}

// WrapperStreamServer creates a single interceptor out of a chain of many interceptors
// Execution is done in right-to-left order
func WrapperStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h := wrapStream(handler, info, interceptors...)
		return h(srv, ss)
	}
}

// wrap h with all specified interceptors
func wrapStream(uh grpc.StreamHandler, info *grpc.StreamServerInfo, interceptors ...grpc.StreamServerInterceptor) grpc.StreamHandler {
	for _, i := range interceptors {
		h := func(current grpc.StreamServerInterceptor, next grpc.StreamHandler) grpc.StreamHandler {
			return func(srv interface{}, stream grpc.ServerStream) error {
				return current(srv, stream, info, next)
			}
		}
		uh = h(i, uh)
	}
	return uh
}

func loggerStreamServerInterceptor(s *Server) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		e := eidFromIncomingContext(ctx)
		// sets new logger instance with eventID
		sublogger := s.logger.With().Str("eid", e).Logger()
		sublogger.Info().Str("method", info.FullMethod).
			Msg(fmt.Sprintf("%s request %s", s.Name(), info.FullMethod))
		// also nice to have a logger available in context
		ctx = sublogger.WithContext(ctx)
		// wrap context
		wrapped := WrapServerStreamWithContext(ss)
		wrapped.SetContext(ctx)
		return handler(srv, wrapped)
	}
}
