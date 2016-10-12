package client_test

import (
	"assert"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/server"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	// Create pluto server
	s := server.NewServer(
		server.Name("client_test_gopher"),
		server.Addr(":65061"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))

	if !testing.Short() {
		// Run Server
		go func() {
			if err := s.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Millisecond * 100)
	}
	result := m.Run()
	if !testing.Short() {
		// Stop Server
		s.Stop()
		time.Sleep(time.Millisecond * 100)
	}
	os.Exit(result)
}

func TestClient(t *testing.T) {

	// Create a grpc client
	c := client.NewClient(
		client.Name("client_test_gopher"),
		client.Description("gopher super client"),
		client.Target("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)
	defer c.Close()

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "client_client_test_gopher", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)

	// Connect
	if err := c.Dial(); err != nil {
		log.Fatal(err)
	}
	// Make a Call
	r, err := c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "Hello client_client_test_gopher", r.Message)
}

func TestHealth(t *testing.T) {
	c := client.NewClient(
		client.Target("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)
	defer c.Close()
	err := c.Dial()
	if err != nil {
		t.Fatal(err)
	}
	h := c.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}
