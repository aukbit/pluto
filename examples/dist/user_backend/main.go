package main

import (
	"flag"
	"log"

	"github.com/gocql/gocql"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/datastore"
	pb "github.com/aukbit/pluto/examples/dist/user_backend/proto"
	"github.com/aukbit/pluto/examples/dist/user_backend/views"
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var (
	grpcPort string
	db       string
	keyspace string
	name     string
)

func init() {
	flag.StringVar(&grpcPort, "grpc_port", ":65060", "grpc listening port")
	flag.StringVar(&db, "db", "cassandra", "datastore service instance")
	flag.StringVar(&keyspace, "keyspace", "pluto_user_backend", "datastore keyspace")
	flag.StringVar(&name, "name", "user_backend", "service name instance")
	flag.Parse()
}

func main() {
	// run service
	if err := service(); err != nil {
		log.Fatal(err)
	}
}

func service() error {

	// Define Pluto Server
	srv := server.New(
		server.Name(name),
		server.Addr(grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &views.UserViews{})
		}))

	cfg := gocql.NewCluster(db)
	cfg.Keyspace = keyspace
	cfg.ProtoVersion = 3
	// Define db connection
	db := datastore.New(
		datastore.Name(name),
		datastore.Cassandra(cfg),
	)

	// Define Pluto Service
	s := pluto.New(
		pluto.Name(name),
		pluto.Description("User backend service is responsible for persist data"),
		pluto.Datastore(db),
		pluto.Servers(srv),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
