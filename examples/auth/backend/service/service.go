package backend

import (
	"crypto/rsa"
	"flag"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/auth/jwt"
	pba "github.com/aukbit/pluto/v6/auth/proto"
	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/examples/auth/backend/views"
	pbu "github.com/aukbit/pluto/v6/examples/user/proto"
	"github.com/aukbit/pluto/v6/server"
	"google.golang.org/grpc"
)

var (
	userTarget string
	grpcPort   string
)

func init() {
	flag.StringVar(&userTarget, "user_target", "127.0.0.1:65080", "user backend address")
	flag.StringVar(&grpcPort, "grpc_port", ":65081", "grpc listening port")
	flag.Parse()
}

// Run runs auth backend service
func Run(pub *rsa.PublicKey, prv *rsa.PrivateKey) error {

	// Define user Client
	clt := client.New(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pbu.NewUserServiceClient(cc)
		}),
		client.Target(userTarget),
	)

	// Define Pluto Server
	srv := server.New(
		server.Addr(grpcPort),
		server.GRPCRegister(func(g *grpc.Server) {
			pba.RegisterAuthServiceServer(g, &backend.AuthViews{})
		}),
		server.UnaryServerInterceptors(jwt.RsaUnaryServerInterceptor(pub, prv)),
	)
	// Logger
	// Define Pluto Service
	s := pluto.New(
		pluto.Name("auth_backend"),
		pluto.Description("Backend service issuing access tokens to the client after successfully authenticating the resource owner and obtaining authorization"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.HealthAddr(":9092"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
