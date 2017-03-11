package backend

import (
	"flag"

	"go.uber.org/zap"

	"github.com/aukbit/pluto"
	pba "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/examples/auth/backend/views"
	pbu "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc"
)

var (
	userTarget = flag.String("user_target", "127.0.0.1:65080", "user backend address")
	grpcPort   = flag.String("grpc_port", ":65081", "grpc listening port")
)

// Run runs auth backend service
func Run() error {
	flag.Parse()

	// Define user Client
	clt := client.New(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pbu.NewUserServiceClient(cc)
		}),
		client.Target(*userTarget),
	)

	// Define Pluto Server
	srv := server.New(
		server.Addr(*grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pba.RegisterAuthServiceServer(g, &backend.AuthViews{})
		}),
	)
	// Logger
	logger, _ := zap.NewDevelopment()
	// Define Pluto Service
	s := pluto.New(
		pluto.Name("auth_backend"),
		pluto.Description("Backend service issuing access tokens to the client after successfully authenticating the resource owner and obtaining authorization"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.Logger(logger),
		pluto.HealthAddr(":9092"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
