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

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible to persist data"),
		pluto.Datastore(*db_addr),
	)

	// 2. Define datastore
	//db := server.NewDatastore()
	// 2. Define grpc Server
	grpcSrv := server.NewGRPCServer(server.Addr(*grpc_port),
		server.RegisterServerFunc(func(srv *grpc.Server){
			pb.RegisterUserServiceServer(srv, &backend.User{Cluster: s.Config().Datastore})
		}),
	)
	// 5. Init service
	s.Init(pluto.Servers(grpcSrv),)

	// 6. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil

}




