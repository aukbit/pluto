package frontend

import (
	"errors"
	"net/http"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/reply"
)

var (
	errBasicAuth              = errors.New("invalid basic authorization header")
	errClientAuthNotAvailable = errors.New("client auth not available")
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	s := ctx.Value("pluto").(pluto.Service)
	l := ctx.Value("logger").(zap.Logger)
	// get authentication from Authorization Header
	u, p, ok := r.BasicAuth()
	if !ok {
		l.Error(errBasicAuth.Error())
		reply.Json(w, r, http.StatusUnauthorized, errBasicAuth)
		return
	}
	// credentials
	cred := &pba.Credentials{Email: u, Password: p}
	// get pluto service from context
	// get gRPC client from service
	c, ok := s.Client("client_auth")
	if !ok {
		l.Error(errClientAuthNotAvailable.Error())
		reply.Json(w, r, http.StatusInternalServerError, errClientAuthNotAvailable)
		return
	}
	// make a call to the backend service
	token, err := c.Call().(pba.AuthServiceClient).Authenticate(ctx, cred)
	if err != nil {
		l.Error(err.Error())
		reply.Json(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, token)
}
