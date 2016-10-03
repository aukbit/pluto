package frontend

import (
	"net/http"

	"bitbucket.org/aukbit/pluto"
	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
	"bitbucket.org/aukbit/pluto/reply"
	"github.com/golang/protobuf/jsonpb"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// credentials
	cred := &pb.Credentials{}
	if err := jsonpb.Unmarshal(r.Body, cred); err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// get pluto service from context
	ctx := r.Context()
	s := ctx.Value("pluto")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_auth")
	// make a call the backend service
	token, err := c.Call().(pb.AuthServiceClient).Authenticate(ctx, cred)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusOK, token)
}
