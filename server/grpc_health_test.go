package server

import (
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "bitbucket.org/aukbit/pluto/server/proto"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v" + in.Name)}, nil
}

func TestHealthGRPC(t *testing.T) {
	s := NewServer(Name("sunshine"), Addr(":65059"),
		GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))
	go func(s Server) {
		//2. Run server
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}

	}(s)
	time.Sleep(time.Millisecond * 100)
	h := s.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}
