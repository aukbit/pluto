package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func Home(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}

func Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reply.Json(w, r, http.StatusOK, fmt.Sprintf("Hello Room %s", ctx.Value("id").(string)))
}

type greeter struct {
	cfg *server.Config
}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("%v: Hello "+in.Name, s.cfg.Name)}, nil
}

func TestServer(t *testing.T) {

	// HTTP server

	// Define Router
	mux := router.NewRouter()
	mux.GET("/home", Home)
	mux.GET("/home/:id", Detail)

	// Create pluto server
	s := server.NewServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":8080"),
		server.Mux(mux),
	)

	cfg := s.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "server_gopher", cfg.Name)
	assert.Equal(t, "gopher super server", cfg.Description)
	assert.Equal(t, ":8080", cfg.Addr)

	// Run server
	go func() {
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	// defer s.Stop()

	// GRPC server
	// Define gRPC server and register
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &greeter{})

	// Create pluto server
	g := server.NewServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":65058"),
		server.GRPCServer(grpcServer),
	)

	// Run Server
	go func() {
		//2. Run server
		if err := g.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	// defer g.Stop()

	// Test
	const URL = "http://localhost:8080"
	var tests = []struct {
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
			Path:         "/home",
			BodyContains: `Hello World`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/home/123",
			BodyContains: `Hello Room 123`,
			Status:       http.StatusOK,
		},
	}
	for _, test := range tests {

		r, err := http.Get(URL + test.Path)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()

		var message string
		if err := json.Unmarshal(b, &message); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, r.StatusCode)
		assert.Equal(t, test.BodyContains, message)

	}
	// Stop server
	// syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

}
