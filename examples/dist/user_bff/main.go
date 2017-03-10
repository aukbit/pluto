package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/discovery"
	pb "github.com/aukbit/pluto/examples/dist/user_bff/proto"
	"github.com/aukbit/pluto/examples/dist/user_bff/views"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
)

var httpPort = flag.String("http_port", ":8080", "backend for frontend http port")
var name = flag.String("name", "user_bff", "service name instance")
var targetName = flag.String("target_name", "user_backend", "target server name instance")
var consulAddr = flag.String("consul_addr", "192.168.99.100:8500", "consul agent address")

func main() {
	flag.Parse()
	// run frontend service
	if err := service(); err != nil {
		log.Fatal(err)
	}
}

func service() error {

	// Define handlers
	mux := router.NewMux()
	mux.GET("/user", views.GetHandler)
	mux.POST("/user", views.PostHandler)
	mux.GET("/user/:id", views.GetHandlerDetail)
	mux.PUT("/user/:id", views.PutHandler)
	mux.DELETE("/user/:id", views.DeleteHandler)

	// Define http server
	srv := server.NewServer(
		server.Name(*name),
		server.Addr(*httpPort),
		server.Mux(mux))

	// Define grpc Client
	clt := client.New(
		client.Name(*name),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.TargetName(*targetName))

	// Define consul
	dis := discovery.NewDiscovery(discovery.Addr(*consulAddr))

	// Define Pluto service
	s := pluto.New(
		pluto.Name(*name),
		pluto.Description("User backend for frontend service is responsible to parse all json data from http requests"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.Discovery(dis),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
