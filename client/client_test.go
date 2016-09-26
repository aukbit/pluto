package client_test

import (
	"testing"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"bitbucket.org/aukbit/pluto/client"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"bitbucket.org/aukbit/pluto/server"
	"reflect"
	"fmt"
)

type greeter struct{
	cfg 			*server.Config
}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("%v: Hello " + in.Name, s.cfg.Name)}, nil
}

func TestClient(t *testing.T){

	// Create a grpc server
	s := server.NewGRPCServer(server.Addr(":65057"))

	s.Init(server.RegisterServerFunc(func (g *grpc.Server){
			pb.RegisterGreeterServer(g, &greeter{cfg: s.Config()})
	}))
	// Run Server
	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// Create a grpc client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
		client.Target("localhost:65057"),
	)
	assert.Equal(t, reflect.TypeOf(client.DefaultClient), reflect.TypeOf(c))

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "client_gopher", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)

	// Connect
	_, err := c.Dial()
	if err != nil {
		log.Fatal(err)
	}
	r, err := c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	//r, err := i.(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "server_default: Hello client_gopher", r.Message)

}