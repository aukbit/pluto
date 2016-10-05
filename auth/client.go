package auth

import (
	pb "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/client"
	"google.golang.org/grpc"
)

// NewClientAuth creates a new default client instance
// to connect to the authorization grpc server
func NewClientAuth(target string) client.Client {
	return client.NewClient(
		client.Name("auth"),
		client.Description("General client to connect to the authorization grpc server"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewAuthServiceClient(cc)
		}),
		client.Target(target))
}
