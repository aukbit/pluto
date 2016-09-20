package backend

import (
	"golang.org/x/net/context"
	"github.com/google/uuid"
	"log"
	"crypto/sha256"
	"encoding/hex"
	"bitbucket.org/aukbit/pluto/datastore"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
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
	defer s.Cluster.Close()
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
	defer s.Cluster.Close()
	// user object
	u := &pb.User{}
	// get data
	if err := s.Cluster.Session().Query(`SELECT id, name, email FROM users WHERE id = ?`, nu.Id).Scan(&u.Id, &u.Name, &u.Email); err != nil {
		log.Fatalf("ERROR ReadUser Query() %v", err)
		return &pb.User{}, err
	}
	return u, nil
}
// ReadUser implements UserServiceServer
func (s *User) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Fatalf("ERROR UpdateUser RefreshSession() %v", err)
		return &pb.User{}, err
	}
	defer s.Cluster.Close()
	// update data
	if err := s.Cluster.Session().Query(`UPDATE users SET name = ?, email = ? WHERE id = ?`, nu.Name, nu.Email, nu.Id).Exec(); err != nil {
		log.Fatalf("ERROR UpdateUser Query() %v", err)
		return &pb.User{}, err
	}
	return nu, nil
}
// DeleteUser implements UserServiceServer
func (s *User) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Fatalf("ERROR DeleteUser RefreshSession() %v", err)
		return &pb.User{}, err
	}
	defer s.Cluster.Close()
	// delete data
	if err := s.Cluster.Session().Query(`DELETE FROM users WHERE id = ?`, nu.Id).Exec(); err != nil {
		log.Fatalf("ERROR DeleteUser Query() %v", err)
		return &pb.User{}, err
	}
	return &pb.User{}, nil
}

// GetUsers implements UserServiceServer
func (s *User) FilterUsers(ctx context.Context, f *pb.Filter) (*pb.Users, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Fatalf("ERROR FilterUsers RefreshSession() %v", err)
		return &pb.Users{}, err
	}
	defer s.Cluster.Close()
	// filter users
	iter := s.Cluster.Session().Query(`SELECT id, name, email FROM users WHERE name = ? ALLOW FILTERING;`, f.Name).Iter()

	users := &pb.Users{}
	u := &pb.User{}
	for iter.Scan(&u.Id, &u.Name, &u.Email) {
		users.Data = append(users.Data, u)
    	}
    	if err := iter.Close(); err != nil {
		log.Fatalf("ERROR FilterUsers Close() %v", err)
		return &pb.Users{}, err
    	}

	return users, nil
}