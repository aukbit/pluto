package jwt

import (
	"context"
	"crypto/rsa"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

var (
	// PublicKeyContextKey is a context key. It can be used in HTTP / GRPC
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of type *rsa.PublicKey.
	PublicKeyContextKey = &contextKey{"public-key"}

	// PrivateKeyContextKey is a context key. It can be used in HTTP / GRPC
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of type *rsa.PublicKey.
	PrivateKeyContextKey = &contextKey{"private-key"}
)

// PublicKeyFromContext retuns public key pointer from a context
func PublicKeyFromContext(ctx context.Context) *rsa.PublicKey {
	return ctx.Value(PublicKeyContextKey).(*rsa.PublicKey)
}

// PrivateKeyFromContext returns private key pointer from a context
func PrivateKeyFromContext(ctx context.Context) *rsa.PrivateKey {
	return ctx.Value(PrivateKeyContextKey).(*rsa.PrivateKey)
}
