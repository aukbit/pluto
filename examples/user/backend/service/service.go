package backend

import (
	"flag"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/examples/user/backend/views"
	pb "github.com/aukbit/pluto/v6/examples/user/proto"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/ext"
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
	// db connection
	cfg := gocql.NewCluster(dbAddr)
	cfg.Keyspace = "examples_user_backend"
	cfg.ProtoVersion = 3

	// Define Pluto Server
	srv := server.New(
		server.Addr(grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterUserServiceServer(g, &backend.UserViews{})
		}),
		server.UnaryServerInterceptors(ext.CassandraUnaryServerInterceptor("cassandra", cfg)),
		server.StreamServerInterceptors(ext.CassandraStreamServerInterceptor("cassandra", cfg)),
	)

	// Define Pluto Service
	s := pluto.New(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible for persist data"),
		pluto.Servers(srv),
		pluto.HealthAddr(":9096"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
