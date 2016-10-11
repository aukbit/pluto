package client_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/server"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	if !testing.Short() {
		// Create pluto server
		s := server.NewServer(
			server.Name("backend"),
			server.Addr(":65061"),
			server.GRPCRegister(func(g *grpc.Server) {
				pb.RegisterGreeterServer(g, &greeter{})
			}))
		// Run Server
		go func() {
			if err := s.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		defer s.Stop()
	}
	result := m.Run()
	if !testing.Short() {
	}
	os.Exit(result)
}

func TestClient(t *testing.T) {

	// Create a grpc client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
		client.Target("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "client_gopher", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)
	//
	// Connect
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	r, err := c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "Hello client_gopher", r.Message)
}

func TestClientDiscover(t *testing.T) {

	// Create a grpc client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
		client.TargetDiscovery("backend_server"), // connect to the server via service discovery
		// client.Target("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "client_gopher", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)
	//
	// Connect
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	r, err := c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "Hello client_gopher", r.Message)
}
