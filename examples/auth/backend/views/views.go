package backend

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"net/http"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/auth/jwt"
	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
	"bitbucket.org/aukbit/pluto/reply"
	"golang.org/x/net/context"
)

type Auth struct {
}

// Authenticate implements authentication
func (s *Auth) Authenticate(ctx context.Context, cre *pb.Credentials) (*pb.Token, error) {
	log.Printf("Login %v", cre.Email)
	// TODO
	// get service from context by service name
	s := ctx.Value("pluto")
	// get gRPC client from service
	c := s.(pluto.Service).Client("client_user")
	// make a call the backend service
	user, err := c.Call().(pb.UserServiceClient).CreateUser(ctx, newUser)
	if err != nil {
		reply.Json(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	reply.Json(w, r, http.StatusCreated, user)

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	token, err := jwt.NewToken(cre.Email, pk)
	if err != nil {
		panic(err)
	}
	return &pb.Token{Jwt: token}, nil
}
