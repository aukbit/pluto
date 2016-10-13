package main

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/common"
	pb "bitbucket.org/aukbit/pluto/examples/dist/user_backend/proto"
	"bitbucket.org/aukbit/pluto/examples/dist/user_backend/views"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var grpcPort = flag.String("grpc_port", ":65060", "grpc listening port")
var db = flag.String("db", "cassandra", "datastore service instance")
var name = flag.String("name", "user_backend", "service name instance")

func main() {
	flag.Parse()
	common.IPaddress()
	common.IP2()
	// run service
	// if err := service(); err != nil {
	// 	log.Fatal(err)
	// }
}

func service() error {
	// Define Pluto Server
	srv := server.NewServer(
		server.Name(*name),
		server.Addr(*grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &views.UserViews{})
		}))

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name(*name),
		pluto.Description("User backend service is responsible for persist data"),
		pluto.DatastoreDiscovery(*db),
		pluto.Servers(srv))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
