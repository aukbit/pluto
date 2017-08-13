package frontend

import (
	"errors"
	"net/http"
	"time"

	"github.com/aukbit/pluto"
	pba "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/reply"
)

var (
	errBasicAuth              = errors.New("invalid basic authorization header")
	errClientAuthNotAvailable = errors.New("client auth not available")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// log := ctx.Value(server.Key("logger")).(*zap.Logger)
	// get authentication from Authorization Header
	u, p, ok := r.BasicAuth()
	if !ok {
		// log.Error(errBasicAuth.Error())
		reply.Json(w, r, http.StatusUnauthorized, errBasicAuth)
		return
	}
	// credentials
	cred := &pba.Credentials{Email: u, Password: p}
	// get pluto service from context
	// get gRPC client from service
	c, ok := pluto.FromContext(ctx).Client("auth")
	if !ok {
		// log.Error(errClientAuthNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientAuthNotAvailable)
		return
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		// log.Error(err.Error())
		reply.Json(w, r, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	// make a call to the backend service
	token, err := c.Stub(conn).(pba.AuthServiceClient).Authenticate(ctx, cred)
	if err != nil {
		// log.Error(err.Error())
		reply.Json(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, token)
}
