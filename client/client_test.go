package client_test

import (
	"testing"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"bitbucket.org/aukbit/pluto/client"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	//"bitbucket.org/aukbit/pluto/server"
	"fmt"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("SayHello")
	return &pb.HelloReply{Message: fmt.Sprintf("%v: Hello " + in.Name)}, nil
}

func TestClient(t *testing.T){

	// Create a grpc server
	// Define gRPC server and register
	//grpcServer := grpc.NewServer()
	//pb.RegisterGreeterServer(grpcServer, &greeter{})
	//// Create pluto server
	//s := server.NewServer(
	//	server.Addr(":65060"),
	//	server.GRPCServer(grpcServer),
	//)
	//// Run Server
	//go func() {
	//	if err := s.Run(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()
	//defer s.Stop()

	// Create a grpc client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
		client.Target("localhost:65060"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "client_gopher", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)
	//
	// Connect
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{})
	//assert.Equal(t, "server_default: Hello client_gopher", r.Message)
}