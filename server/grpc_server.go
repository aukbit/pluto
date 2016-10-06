package server

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/uber-go/zap"
)

type B struct{}

// serveGRPC serves *grpc.Server
func (ds *defaultServer) serveGRPC(ln net.Listener) (err error) {

	// register
	//
	srv := grpc.NewServer(grpc.UnaryInterceptor(WrapUnaryInterceptor("logger", ds.logger)))
	// srv := grpc.NewServer()
	ds.cfg.GRPCRegister(srv)

	// srv := ds.cfg.GRPCServer
	// add go routine to WaitGroup
	ds.wg.Add(1)
	go func() {
		defer ds.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
			return
		}
	}()
	return nil
}

// WrapUnaryInterceptor
func WrapUnaryInterceptor(key interface{}, val interface{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, key, val)
		l := ctx.Value("logger")
		md, ok := metadata.FromContext(ctx)
		if ok {
			e, ok := md["event"]
			if ok {
				l.(zap.Logger).Info("request", zap.String("event", e[0]))
			}
		}
		return handler(ctx, req)
	}
}

// func AuthUnaryInterceptor(
// 	ctx context.Context,
// 	req interface{},
// 	info *grpc.UnaryServerInfo,
// 	handler grpc.UnaryHandler,
// ) (interface{}, error) {
//
// 	// retrieve metadata from context
// 	md, ok := metadata.FromContext(ctx)
//
// 	// validate 'authorization' metadata
// 	// like headers, the value is an slice []string
// 	uid, err := MyValidationFunc(md["authorization"])
// 	if err != nil {
// 		return nil, grpc.Errorf(codes.Unauthenticated, "authentication required")
// 	}
//
// 	// add user ID to the context
// 	newCtx := context.WithValue(ctx, "user_id", uid)
//
// 	// handle scopes?
// 	// ...
// 	return handler(newCtx, req)
// }
