package client_test

import (
	"testing"
	"pluto/client"
	"github.com/paulormart/assert"
	pb "pluto/server/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
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
	assert.Equal(t, "gopher.client.grpc", cfg.Name)
	assert.Equal(t, "gopher super client", cfg.Description)

	//2.
	i, err := c.Dial()
	if err != nil {
		log.Fatal(err)
	}
	r, err := i.(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: cfg.Name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r.Message)

}