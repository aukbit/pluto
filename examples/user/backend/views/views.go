package backend

import (
	"golang.org/x/net/context"
	pb "pluto/examples/user/proto"
	"log"
)


type User struct{
}

// CreateUser implements UserServiceServer
func (s *User) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	//serv := ctx.Value("frontend.pluto")
	log.Printf("CreateUser %v", ctx)
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// ReadUser implements UserServiceServer
func (s *User) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// ReadUser implements UserServiceServer
func (s *User) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// DeleteUser implements UserServiceServer
func (s *User) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{}, nil
}
