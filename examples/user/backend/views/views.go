package backend

import (
	"golang.org/x/net/context"
	"github.com/google/uuid"
	pb "pluto/examples/user/proto"
	"pluto/datastore"
	"log"
	"crypto/sha256"
	"encoding/hex"
)


type User struct {
	Cluster		datastore.Datastore
}

// CreateUser implements UserServiceServer
func (s *User) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Fatalf("ERROR CreateUser RefreshSession() %v", err)
		return &pb.User{}, err
	}
	// generate user id uuid
	newId := uuid.New().String()
	// hash password
	h := sha256.New()
	h.Write([]byte(nu.Password))
	sha256_hash := hex.EncodeToString(h.Sum(nil))
	// persist data
	if err := s.Cluster.Session().Query(`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		newId, nu.Name, nu.Email, sha256_hash).Exec(); err != nil {
		log.Fatalf("ERROR CreateUser Query() %v", err)
		return &pb.User{}, err
	    }
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: newId}, nil
}
// ReadUser implements UserServiceServer
func (s *User) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Fatalf("ERROR ReadUser RefreshSession() %v", err)
		return &pb.User{}, err
	}
	u := &pb.User{}
	// get data
	if err := s.Cluster.Session().Query(`SELECT id, name, email FROM users WHERE id = ?`, nu.Id).Scan(u.Id, u.Name, u.Email); err != nil {
		log.Fatalf("ERROR ReadUser Query() %v", err)
		return &pb.User{}, err
	}
	return u, nil
}
// ReadUser implements UserServiceServer
func (s *User) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: "123"}, nil
}
// DeleteUser implements UserServiceServer
func (s *User) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	return &pb.User{}, nil
}
