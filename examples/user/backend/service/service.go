package backend

import (
	"flag"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/examples/user/backend/views"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/server"
	"github.com/gocql/gocql"
	"google.golang.org/grpc"
)

var (
	dbAddr   string
	grpcPort string
)

func init() {
	flag.StringVar(&dbAddr, "db_addr", "127.0.0.1", "cassandra address")
	flag.StringVar(&grpcPort, "grpc_port", ":65087", "grpc listening port")
	flag.Parse()
}

func Run() error {
	// Define Pluto Server
	srv := server.New(
		server.Addr(grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &backend.UserViews{})
		}),
	)
	// db connection
	cfg := gocql.NewCluster(dbAddr)
	cfg.Keyspace = "examples_user_backend"
	cfg.ProtoVersion = 3
	db := datastore.New(
		datastore.Cassandra(cfg),
	)
	// logger
	// logger, _ := zap.NewDevelopment()
	// Define Pluto Service
	s := pluto.New(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible for persist data"),
		pluto.Datastore(db),
		pluto.Servers(srv),
		// pluto.Logger(logger),
		pluto.HealthAddr(":9096"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
