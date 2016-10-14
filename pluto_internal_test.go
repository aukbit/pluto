package pluto

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

	"github.com/paulormart/assert"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/discovery"
	"bitbucket.org/aukbit/pluto/server"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	// Create pluto server
	srvHTTP := server.NewServer(
		server.Name("http"),
		server.Description("gopher super server"),
		server.Addr(":8080"),
	)
	// Create grpc pluto server
	srvGRPC := server.NewServer(
		server.Name("grpc"),
		server.Description("grpc super server"),
		server.Addr(":65060"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))
	// Create grpc pluto client
	clt := client.NewClient(
		client.Name("grpc"),
		// client.Target("localhost:65060"),
		client.TargetName("grpc"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewGreeterClient(cc)
		}),
	)
	// Define service Discovery
	d := discovery.NewDiscovery(discovery.Addr("192.168.99.100:8500"))
	// Define Pluto Service
	s := NewService(
		Name("gopher"),
		Servers(srvHTTP),
		Servers(srvGRPC),
		Clients(clt),
		Discovery(d),
	)

	if !testing.Short() {
		// Run Server
		go func() {
			if err := s.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Second)
	}
	result := m.Run()
	if !testing.Short() {
		// Stop Server
		s.Stop()
		time.Sleep(time.Millisecond * 100)
	}
	os.Exit(result)
}

const URL = "http://localhost:9090/_health"

func TestHealth(t *testing.T) {

	var tests = []struct {
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
			Path:         "/server/server_health",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/server_http",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/server_grpc",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/server/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
		},
		{
			Path:         "/client/client_grpc",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/client/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
		},
		{
			Path:         "/pluto/pluto_gopher",
			BodyContains: `SERVING`,
			Status:       http.StatusOK,
		},
		{
			Path:         "/pluto/something_wrong",
			BodyContains: `UNKNOWN`,
			Status:       http.StatusNotFound,
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
		hcr := &healthpb.HealthCheckResponse{}
		if err := json.Unmarshal(b, hcr); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, r.StatusCode)
		assert.Equal(t, test.BodyContains, hcr.Status.String())
	}
}
