package frontend

import (
	"errors"
	"net/http"

	"bitbucket.org/aukbit/pluto"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/reply"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// get authentication from Authorization Header
	u, p, ok := r.BasicAuth()
	if !ok {
		reply.Json(w, r, http.StatusUnauthorized, `{}`)
		return
	}
	// credentials
	cred := &pba.Credentials{Email: u, Password: p}
	// get pluto service from context
	ctx := r.Context()
	s := ctx.Value("pluto")
	// get gRPC client from service
	c, ok := s.(pluto.Service).Client("client_auth")
	if !ok {
		err := errors.New("Client auth not available")
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// make a call to the backend service
	token, err := c.Call().(pba.AuthServiceClient).Authenticate(ctx, cred)
	if err != nil {
		reply.Json(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, token)
}
