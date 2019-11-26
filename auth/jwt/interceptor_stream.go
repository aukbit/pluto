package jwt

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BearerTokenStreamServerInterceptor makes bearer token available in grpc context
func BearerTokenStreamServerInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return streamer(ctx, desc, cc, method, opts...)
		}
		auth, ok := md["authorization"]
		if !ok {
			return streamer(ctx, desc, cc, method, opts...)
		}
		t, ok := parseBearerAuth(auth[0])
		if !ok {
			return streamer(ctx, desc, cc, method, opts...)
		}
		ctx = context.WithValue(ctx, TokenContextKey, t)

		return streamer(ctx, desc, cc, method, opts...)
	}
}
