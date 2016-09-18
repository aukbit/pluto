package backend

import (
	"pluto"
	"pluto/server"
	pb "pluto/examples/user/proto"
	"google.golang.org/grpc"
	"pluto/examples/user/backend/views"
)


func Run() error {

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible to persist data"),
		pluto.Datastore("127.0.0.1"),
	)

	// 2. Define datastore
	//db := server.NewDatastore()
	// 2. Define grpc Server
	grpcSrv := server.NewGRPCServer(server.Addr(":65060"),
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




