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

	// TokenContextKey is a context key. It can be used in HTTP / GRPC
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of type string.
	TokenContextKey = &contextKey{"token-key"}
)

// PublicKeyFromContext retuns public key pointer from a context if it exists.
func PublicKeyFromContext(ctx context.Context) (p *rsa.PublicKey, ok bool) {
	p, ok = ctx.Value(PublicKeyContextKey).(*rsa.PublicKey)
	return
}

// PrivateKeyFromContext returns private key pointer from a context
func PrivateKeyFromContext(ctx context.Context) (p *rsa.PrivateKey, ok bool) {
	p, ok = ctx.Value(PrivateKeyContextKey).(*rsa.PrivateKey)
	return
}

// TokenFromContext retuns token from a context if it exists.
func TokenFromContext(ctx context.Context) (p string, ok bool) {
	p, ok = ctx.Value(TokenContextKey).(string)
	return
}
