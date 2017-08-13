package pluto

import (
	"context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

var (
	// PlutoContextKey is a context key. It can be used in HTTP / GRPC
	// handlers with context.WithValue to access the server that
	// started the handler. The associated value will be of type *Server.
	PlutoContextKey = &contextKey{"pluto-server"}
)

// FromContext returns pluto service pointer from a context
func FromContext(ctx context.Context) *Service {
	return ctx.Value(PlutoContextKey).(*Service)
}

// WithContext returns a copy of ctx with pluto service associated.
func (s *Service) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, PlutoContextKey, s)
}
