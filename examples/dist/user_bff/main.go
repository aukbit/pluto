package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/client"
	pb "github.com/aukbit/pluto/examples/dist/user_bff/proto"
	"github.com/aukbit/pluto/examples/dist/user_bff/views"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
)

var (
	httpPort   string
	name       string
	target     string
	consulAddr string
)

func init() {
	flag.StringVar(&httpPort, "http_port", ":8080", "backend for frontend http port")
	flag.StringVar(&name, "name", "user_bff", "service name instance")
	flag.StringVar(&target, "target_name", "user_backend:65060", "target server name instance")
	flag.StringVar(&consulAddr, "consul_addr", "192.168.99.100:8500", "consul agent address")
	flag.Parse()
}

func main() {
	// run frontend service
	if err := service(); err != nil {
		log.Fatal(err)
	}
}

func service() error {

	// Define handlers
	mux := router.New()
	mux.GET("/user", views.GetHandler)
	mux.POST("/user", views.PostHandler)
	mux.GET("/user/:id", views.GetHandlerDetail)
	mux.PUT("/user/:id", views.PutHandler)
	mux.DELETE("/user/:id", views.DeleteHandler)

	// Define http server
	srv := server.New(
		server.Name(name),
		server.Addr(httpPort),
		server.Mux(mux))

	// Define grpc Client
	clt := client.New(
		client.Name(name),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target(target))

	// Define Pluto service
	s := pluto.New(
		pluto.Name(name),
		pluto.Description("User backend for frontend service is responsible to parse all json data from http requests"),
		pluto.Servers(srv),
		pluto.Clients(clt),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
