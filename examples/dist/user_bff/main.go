package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/client"
	pb "bitbucket.org/aukbit/pluto/examples/dist/user_bff/proto"
	"bitbucket.org/aukbit/pluto/examples/dist/user_bff/views"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

var http_port = flag.String("http_port", ":8080", "backend for frontend http port")

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

	// define http server
	srv := server.NewServer(
		server.Name("user_bff"),
		server.Addr(*http_port),
		server.Mux(mux))

	// Define grpc Client
	clt := client.NewClient(
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.TargetDiscovery("server_user_backend"),
	)
	// Define Pluto service
	s := pluto.NewService(
		pluto.Name("user_bff"),
		pluto.Description("User backend for frontend service is responsible to parse all json data from http requests"),
		pluto.Servers(srv),
		pluto.Clients(clt))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
