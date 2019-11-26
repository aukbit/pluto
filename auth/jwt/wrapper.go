package jwt

import (
	"context"
	"errors"
	"net/http"

	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server/router"
	"google.golang.org/grpc"
)

var (
	errInvalidBearer = errors.New("invalid bearer authorization header")
)

// WrapBearerToken adds token to the context.
func WrapBearerToken(h router.HandlerFunc) router.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get jwt token from Authorization header
		t, ok := BearerAuth(r)
		if !ok {
			reply.Json(w, r, http.StatusUnauthorized, errInvalidBearer)
			return
		}
		ctx := context.WithValue(r.Context(), TokenContextKey, t)

		// pass execution to the original handler
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

// WrapBearerTokenErr adds token to the context.
func WrapBearerTokenErr(h router.WrapErr) router.WrapErr {
	return func(w http.ResponseWriter, r *http.Request) *router.Err {
		// Get jwt token from Authorization header
		t, ok := BearerAuth(r)
		if !ok {
			return &router.Err{
				Err:     errInvalidBearer,
				Status:  http.StatusBadRequest,
				Type:    "authentication_error",
				Message: grpc.ErrorDesc(errInvalidBearer),
			}
		}
		ctx := context.WithValue(r.Context(), TokenContextKey, t)

		// pass execution to the original handler
		h.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}
}
