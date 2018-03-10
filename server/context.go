package server

import (
	"context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

var (
	// ServerContextKey is a context key. It can be used in HTTP / GRPC
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of type *Server.
	ServerContextKey = &contextKey{"pluto-server"}
)

// FromContext returns server instance from a context
func FromContext(ctx context.Context) *Server {
	return ctx.Value(ServerContextKey).(*Server)
}

// WithContext returns a copy of ctx with server associated.
func (s *Server) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ServerContextKey, s)
}

// FromContextAny returns from context the interface value to which the key is
// associated.
func FromContextAny(ctx context.Context, key string) interface{} {
	return ctx.Value(contextKey{key})
}

// WithContextAny returns a copy of parent ctx in which the value associated
// with key is val.
func WithContextAny(ctx context.Context, key string, val interface{}) context.Context {
	return context.WithValue(ctx, contextKey{key}, val)
}
