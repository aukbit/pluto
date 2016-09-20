package client_test

import (
	"testing"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"bitbucket.org/aukbit/pluto/client"
	pb "bitbucket.org/aukbit/pluto/server/proto"
)

func TestClient(t *testing.T){

	//1. create a client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
		client.Target("127.0.0.1:65057"),
	)
	//assert.Equal(t, reflect.TypeOf(client.DefaultClient), reflect.TypeOf(c))

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "gopher.client", cfg.Name)
	assert.Equal(t, "grpc", cfg.Format)
	assert.Equal(t, "gopher super client", cfg.Description)

	//2.
	_, err := c.Dial()
	if err != nil {
		log.Fatal(err)
	}
	r, err := c.Call().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	//r, err := i.(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r.Message)

}