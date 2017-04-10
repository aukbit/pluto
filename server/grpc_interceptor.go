package server

import (
	"github.com/aukbit/pluto/common"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// WrapperUnaryServer creates a single interceptor out of a chain of many interceptors
// Execution is done in right-to-left order
func WrapperUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		h := wrap(handler, info, interceptors...)
		return h(ctx, req)
	}
}

// wrap h with all specified interceptors
func wrap(uh grpc.UnaryHandler, info *grpc.UnaryServerInfo, interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryHandler {
	for _, i := range interceptors {
		h := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return current(ctx, req, info, next)
			}
		}
		uh = h(i, uh)
	}
	return uh
}

func loggerUnaryServerInterceptor(srv *Server) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// get or create unique event id for every request
		e, ctx := common.GetOrCreateEventID(ctx)
		// create new log instance with eventID
		l := srv.logger.With(zap.String("event", e))
		l.Info("request", zap.String("method", info.FullMethod))
		// also nice to have a logger available in context
		ctx = context.WithValue(ctx, Key("logger"), l)
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
