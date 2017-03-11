package client_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/server"
	pb "github.com/aukbit/pluto/test/proto"
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
	logger, _ := zap.NewDevelopment()
	// Create pluto server
	s := server.New(
		server.Name("client_test_gopher"),
		server.Addr(":65061"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}),
		server.Logger(logger),
	)

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

	logger, _ := zap.NewDevelopment()
	// Create a grpc client
	c := client.New(
		client.Logger(logger),
		client.Name("client_test_gopher"),
		client.Description("gopher super client"),
		client.Targets("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)
	defer c.Close()

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "client_test_gopher_client", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)

	// Connect
	if err := c.Dial(); err != nil {
		t.Fatal(err)
	}
	// request a conn from client
	conn := c.Request()
	// when finished with request call done on connector
	defer c.Done(conn)
	// assert proto type
	client := conn.Client().(pb.GreeterClient)
	// call a method
	r, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello client_test_gopher_client", r.Message)
	//
	// or call method Call() directly if load balencer is serving only one service
	r, err = c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello client_test_gopher_client", r.Message)
}

func TestHealth(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	c := client.New(
		client.Logger(logger),
		client.Targets("localhost:65061"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)
	defer c.Close()
	// Connect
	err := c.Dial()
	if err != nil {
		t.Fatal(err)
	}
	h := c.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}
