package backend

import (
	"google.golang.org/grpc"
	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/examples/user/backend/views"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"flag"
)

var db_addr = flag.String("db_addr", "127.0.0.1", "cassandra address")
var grpc_port = flag.String("grpc_port", ":65060", "grpc listening port")

func Run() error {
	flag.Parse()

	// GRPC server
	// Define gRPC server and register
	grpcServer := grpc.NewServer()
	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(*grpc_port), server.GRPCServer(grpcServer))
	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible to persist data"),
		pluto.Datastore(*db_addr),
		pluto.Servers(grpcSrv),
	)
	// Register grpc Server
	pb.RegisterUserServiceServer(grpcServer, &backend.User{Cluster: s.Config().Datastore})

	// 5. Init service
	// TODO remove init method redundant
	s.Init()

	// 6. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil

}




