package backend

import (
	"pluto"
	"pluto/server"
	pb "pluto/examples/user/proto"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type user struct{}

// CreateUser implements UserServiceServer
func (s *user) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// ReadUser implements UserServiceServer
func (s *user) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// ReadUser implements UserServiceServer
func (s *user) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// DeleteUser implements UserServiceServer
func (s *user) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{}, nil
}

func Run() error {

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("backend"),
		pluto.Description("Backend service is responsible to persist data"),
	)

	// 2. Define grpc Server
	grpcSrv := server.NewGRPCServer(server.Addr(":65060"),
		server.RegisterServerFunc(func(srv *grpc.Server){
			pb.RegisterUserServiceServer(srv, &user{})
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


