package ext

import (
	context "golang.org/x/net/context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
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
