package frontend

import (
	"flag"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/examples/user/frontend/views"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	"google.golang.org/grpc"
)

var target = flag.String("target", "127.0.0.1:65065", "backend address")
var http_port = flag.String("http_port", ":8087", "frontend http port")

func Run() error {
	flag.Parse()

	// Define handlers
	mux := router.NewMux()
	mux.GET("/user", frontend.GetHandler)
	mux.POST("/user", frontend.PostHandler)
	mux.GET("/user/:id", frontend.GetHandlerDetail)
	mux.PUT("/user/:id", frontend.PutHandler)
	mux.DELETE("/user/:id", frontend.DeleteHandler)

	// define http server
	srv := server.NewServer(
		server.Name("api"),
		server.Addr(*http_port),
		server.Mux(mux),
	)

	// Define grpc Client
	clt := client.NewClient(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target(*target),
	)
	// Define Pluto service
	s := pluto.New(
		pluto.Name("frontend"),
		pluto.Description("Frontend service is responsible to parse all json data to regarding users to internal services"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.Development(),
		pluto.HealthAddr(":9097"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
