package jwt

import (
	"context"
	"crypto/rsa"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RsaUnaryServerInterceptor makes rsa public and private keys available in grpc context
func RsaUnaryServerInterceptor(a *rsa.PublicKey, b *rsa.PrivateKey) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, PublicKeyContextKey, a)
		ctx = context.WithValue(ctx, PrivateKeyContextKey, b)
		return handler(ctx, req)
	}
}

// BearerTokenUnaryServerInterceptor makes bearer token available in grpc context
func BearerTokenUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		auth, ok := md["authorization"]
		if !ok {
			return handler(ctx, req)
		}
		t, ok := parseBearerAuth(auth[0])
		if !ok {
			return handler(ctx, req)
		}
		ctx = context.WithValue(ctx, TokenContextKey, t)

		return handler(ctx, req)
	}
}
