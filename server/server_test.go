package server_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	pb "github.com/aukbit/pluto/test/proto"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const URL = "http://localhost:8085"

func Home(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}

func Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reply.Json(w, r, http.StatusOK, fmt.Sprintf("Hello Room %s", ctx.Value("id").(string)))
}

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	// Define Router
	mux := router.NewMux()
	mux.GET("/home", Home)
	mux.GET("/home/:id", Detail)

	// Create pluto server
	s := server.NewServer(
		server.Name("http"),
		server.Description("gopher super server"),
		server.Addr(":8085"),
		server.Mux(mux),
	)
	// Create grpc pluto server
	g := server.NewServer(
		server.Name("grpc"),
		server.Description("grpc super server"),
		server.Addr(":65050"),
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
		go func() {
			//2. Run server
			if err := g.Run(); err != nil {
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
		g.Stop()
		time.Sleep(time.Millisecond * 100)
	}
	os.Exit(result)
}

func TestHttpHealthCheck(t *testing.T) {
	r, err := http.Get(URL + "/_health")
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	hcr := &healthpb.HealthCheckResponse{}
	if err := json.Unmarshal(b, hcr); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "SERVING", hcr.Status.String())
}

func TestHttpServer(t *testing.T) {
	// Test
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

func TestGrpcHealthCheck(t *testing.T) {
	time.Sleep(time.Second)
	conn, err := grpc.Dial("localhost:65050", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := healthpb.NewHealthClient(conn)

	h, err := c.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "SERVING", h.Status.String())
}
