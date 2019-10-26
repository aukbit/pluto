package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/auth/jwt"
	pba "github.com/aukbit/pluto/v6/auth/proto"
	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server/router"
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
			c, ok := pluto.FromContext(ctx).Client("auth")
			if !ok {
				err := errors.New("Authorization service not available")
				reply.Json(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			// dial
			conn, err := c.Dial(client.Timeout(5 * time.Second))
			if err != nil {
				reply.Json(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			defer conn.Close()
			// make a call to the Auth backend service
			v, err := c.Stub(conn).(pba.AuthServiceClient).Verify(ctx, &pba.Token{Jwt: t})
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
