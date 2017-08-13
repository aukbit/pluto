package router

import (
	"context"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	param string
}

// FromContext returns server instance from a context
func FromContext(ctx context.Context, param string) string {
	return ctx.Value(contextKey{param}).(string)
}
