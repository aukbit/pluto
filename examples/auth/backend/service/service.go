package backend

import (
	"flag"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/examples/auth/backend/views"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var (
	userTarget = flag.String("user_target", "127.0.0.1:65080", "user backend address")
	grpcPort   = flag.String("grpc_port", ":65081", "grpc listening port")
	logger     = zap.New(zap.NewJSONEncoder())
)

// Run runs auth backend service
func Run() error {
	flag.Parse()

	// Define user Client
	userClient := client.NewClient(
		client.Name("user"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pbu.NewUserServiceClient(cc)
		}),
		client.Target(*userTarget))

	// GRPC server
	// Define gRPC server
	grpcServer := grpc.NewServer()

	// Define Pluto Server
	grpcSrv := server.NewServer(server.Addr(*grpcPort), server.GRPCServer(grpcServer))

	// Define Pluto Service
	s := pluto.NewService(
		pluto.Name("auth_backend"),
		pluto.Description("Backend service is responsible verify for persist data"),
		pluto.Servers(grpcSrv),
		pluto.Clients(userClient))

	// Register Auth service passing user backend client

	pba.RegisterAuthServiceServer(grpcServer, &backend.Auth{Clt: userClient})

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
