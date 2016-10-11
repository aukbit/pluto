package main

import (
	"flag"
	"log"

	"bitbucket.org/aukbit/pluto"
	pb "bitbucket.org/aukbit/pluto/examples/dist/user_backend/proto"
	"bitbucket.org/aukbit/pluto/examples/dist/user_backend/views"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var db_addr = flag.String("db_addr", "127.0.0.1", "cassandra address")
var grpc_port = flag.String("grpc_port", ":65060", "grpc listening port")

func main() {
	flag.Parse()

	// run service
	if err := service(); err != nil {
		log.Fatal(err)
	}
}

func service() error {
	// Define Pluto Server
	srv := server.NewServer(
		server.Addr(*grpc_port),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &views.UserViews{})
		}))

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("user_backend"),
		pluto.Description("User backend service is responsible for persist data"),
		pluto.DatastoreDiscovery("cassandra"),
		pluto.Servers(srv))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
