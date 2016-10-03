package backend

import (
	"log"

	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
	"golang.org/x/net/context"
)

type Auth struct {
}

// Authenticate implements authentication
func (s *Auth) Authenticate(ctx context.Context, cre *pb.Credentials) (*pb.Token, error) {
	log.Printf("Login %v", cre.Email)
	return &pb.Token{Token: "123"}, nil
}
