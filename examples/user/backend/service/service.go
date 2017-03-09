package backend

import (
	"flag"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/examples/user/backend/views"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var db_addr = flag.String("db_addr", "127.0.0.1", "cassandra address")
var grpc_port = flag.String("grpc_port", ":65060", "grpc listening port")

func Run() error {
	flag.Parse()

	// Define Pluto Server
	srv := server.NewServer(
		server.Addr(*grpc_port),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &backend.UserViews{})
		}))
	// db connection
	db := datastore.NewDatastore(
		datastore.Target(*db_addr),
		datastore.Keyspace("examples_user_backend"))
	// Define Pluto Service
	s := pluto.New(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible for persist data"),
		pluto.Datastore(db),
		pluto.Servers(srv),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
