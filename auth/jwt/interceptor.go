package jwt

import (
	"crypto/rsa"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// RsaUnaryServerInterceptor makes rsa public and private keys available in grpc context
func RsaUnaryServerInterceptor(a *rsa.PublicKey, b *rsa.PrivateKey) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = context.WithValue(ctx, PublicKeyContextKey, a)
		ctx = context.WithValue(ctx, PrivateKeyContextKey, b)
		return handler(ctx, req)
	}
}
