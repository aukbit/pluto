package backend

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/examples/auth/backend/views"
	pb "bitbucket.org/aukbit/pluto/examples/auth/proto"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var grpcPort = flag.String("grpc_port", ":65081", "grpc listening port")

// Run runs auth backend service
func Run() error {
	flag.Parse()

	// GRPC server
	// Define gRPC server and register
	grpcServer := grpc.NewServer()

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("auth_backend"),
		pluto.Description("Backend service is responsible for persist data"),
	)
	// Register grpc Server
	pb.RegisterAuthServiceServer(grpcServer, &backend.Auth{})

	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(*grpcPort), server.GRPCServer(grpcServer))

	// 5. Init service
	// TODO remove init method redundant
	s.Init(pluto.Servers(grpcSrv))
	// 6. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
