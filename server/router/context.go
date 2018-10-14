package router

import (
	"context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	param string
}

// FromContext returns value for repective param available in context
// NOte: deprecated, please use FromContextParam
func FromContext(ctx context.Context, param string) string {
	return ctx.Value(contextKey{param}).(string)
}

// FromContextParam returns value for repective param available in context
func FromContextParam(ctx context.Context, param string) string {
	if v, ok := ctx.Value(contextKey{param}).(string); ok {
		return v
	}
	return ""
}

// WithContextParam returns a copy of ctx with param value associated.
func WithContextParam(ctx context.Context, param, val string) context.Context {
	return context.WithValue(ctx, contextKey{param}, val)
}
