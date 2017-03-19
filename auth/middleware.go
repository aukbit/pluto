package auth

import (
	"errors"
	"net/http"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/auth/jwt"
	pba "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
)

// MiddlewareBearerAuth Middleware to validate all handlers with
// Authorization: Bearer jwt
func MiddlewareBearerAuth() router.Middleware {
	return func(h router.HandlerFunc) router.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// get jwt token from Authorization header
			t, ok := jwt.BearerAuth(r)
			if !ok {
				err := errors.New("Invalid Bearer Authorization header")
				reply.Json(w, r, http.StatusUnauthorized, err.Error())
				return
			}
			// verify if token is valid with Auth backend service
			ctx := r.Context()
			// get gRPC Auth Client from pluto service context
			c, ok := ctx.Value("pluto").(*pluto.Service).Client("auth")
			if !ok {
				err := errors.New("Authorization service not available")
				reply.Json(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			// make a call to the Auth backend service
			v, err := c.Call().(pba.AuthServiceClient).Verify(ctx, &pba.Token{Jwt: t})
			if err != nil {
				reply.Json(w, r, http.StatusUnauthorized, err.Error())
				return
			}
			if !v.IsValid {
				err := errors.New("Invalid token")
				reply.Json(w, r, http.StatusUnauthorized, err.Error())
				return
			}
			//
			h.ServeHTTP(w, r)
		}
	}
}
