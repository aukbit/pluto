package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

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
	mux := router.NewMux()
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
	defer s.Stop()
	// Create pluto server
	g := server.NewServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":65058"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))

	// Run Server
	go func() {
		//2. Run server
		if err := g.Run(); err != nil {
			log.Fatal(err)
		}

	}()
	defer g.Stop()
	// wait a bit for service discover
	time.Sleep(time.Second)
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
}
