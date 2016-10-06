package frontend

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/examples/user/frontend/views"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"google.golang.org/grpc"
)

var target = flag.String("target", "127.0.0.1:65060", "backend address")
var http_port = flag.String("http_port", ":8080", "frontend http port")

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
		server.Mux(mux))

	// Define grpc Client
	clt := client.NewClient(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target(*target),
	)
	// Define Pluto service
	s := pluto.NewService(
		pluto.Name("frontend"),
		pluto.Description("Frontend service is responsible to parse all json data to regarding users to internal services"),
		pluto.Servers(srv),
		pluto.Clients(clt))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
