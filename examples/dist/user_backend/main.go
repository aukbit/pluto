package main

import (
	"flag"
	"log"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/discovery"
	pb "github.com/aukbit/pluto/examples/dist/user_backend/proto"
	"github.com/aukbit/pluto/examples/dist/user_backend/views"
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var grpcPort = flag.String("grpc_port", ":65060", "grpc listening port")
var db = flag.String("db", "cassandra", "datastore service instance")
var keyspace = flag.String("keyspace", "pluto_user_backend", "datastore keyspace")
var name = flag.String("name", "user_backend", "service name instance")
var consulAddr = flag.String("consul_addr", "192.168.99.100:8500", "consul agent address")

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
		server.Name(*name),
		server.Addr(*grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &views.UserViews{})
		}))

	// Define db connection
	db := datastore.NewDatastore(
		datastore.Name(*name),
		datastore.TargetName(*db),
		datastore.Keyspace(*keyspace))

	// Define consul
	dis := discovery.NewDiscovery(discovery.Addr(*consulAddr))

	// Define Pluto Service
	s := pluto.New(
		pluto.Name(*name),
		pluto.Description("User backend service is responsible for persist data"),
		pluto.Datastore(db),
		pluto.Servers(srv),
		pluto.Discovery(dis),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
