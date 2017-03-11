package server

import (
	"fmt"
	"log"
	"testing"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc"

	"github.com/paulormart/assert"

	pb "github.com/aukbit/pluto/test/proto"

	"golang.org/x/net/context"
)

func TestHealthHTTP(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	s := New(
		Name("awesome"),
		Addr(":8082"),
		Logger(logger),
	)
	go func(s *Server) {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}(s)
	defer s.Stop()
	time.Sleep(time.Second)
	h := s.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v" + in.Name)}, nil
}

func TestHealthGRPC(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	s := New(
		Name("sunshine"),
		Addr(":65059"),
		GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}),
		Logger(logger),
	)
	go func(s *Server) {
		//2. Run server
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}

	}(s)
	time.Sleep(time.Millisecond * 100)
	h := s.Health()
	assert.Equal(t, "SERVING", h.Status.String())
}
