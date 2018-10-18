package jwt

import (
	"context"
	"crypto/rsa"
	"net/http"

	"github.com/aukbit/pluto/server/router"
)

// Middleware adds *rsa.PublicKey and *rsa.PrivateKey to the context.
func Middleware(a *rsa.PublicKey, b *rsa.PrivateKey) router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), PublicKeyContextKey, a)
			ctx = context.WithValue(ctx, PrivateKeyContextKey, b)
			// pass execution to the original handler
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
