package auth

import (
	pb "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/client"
	"google.golang.org/grpc"
)

// NewClientAuth creates a new default client instance
// to connect to the authorization grpc server
func NewClientAuth(target string) client.Client {
	return client.NewClient(
		client.Name("auth"),
		client.Description("General client to connect to the authorization grpc server"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewAuthServiceClient(cc)
		}),
		client.Target(target))
}
