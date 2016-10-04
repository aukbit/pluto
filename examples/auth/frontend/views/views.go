package frontend

import (
	"net/http"

	"bitbucket.org/aukbit/pluto"
	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
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
	cred := &pb.Credentials{Email: u, Password: p}
	// get pluto service from context
	ctx := r.Context()
	s := ctx.Value("pluto")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_auth")
	// make a call to the backend service
	token, err := c.Call().(pb.AuthServiceClient).Authenticate(ctx, cred)
	if err != nil {
		reply.Json(w, r, http.StatusUnauthorized, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, token)
}
