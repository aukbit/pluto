package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"bitbucket.org/aukbit/pluto/datastore"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// User struct
type User struct {
	Cluster datastore.Datastore
}

// CreateUser implements UserServiceServer
func (s *User) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR CreateUser RefreshSession() %v", err)
		return &pb.User{}, err
	}
	defer s.Cluster.Close()
	// generate user id uuid
	newID := uuid.New().String()
	// hash password
	passwordHash := hashPassword(nu.Password)
	// persist data
	if err := s.Cluster.Session().Query(`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		newID, nu.Name, nu.Email, passwordHash).Exec(); err != nil {
		log.Printf("ERROR CreateUser Query() %v", err)
		return &pb.User{}, err
	}
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: newID}, nil
}

// ReadUser implements UserServiceServer
func (s *User) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR ReadUser RefreshSession() %v", err)
		return nu, err
	}
	defer s.Cluster.Close()
	// user object
	u := &pb.User{}
	// get data
	if err := s.Cluster.Session().Query(`SELECT id, name, email FROM users WHERE id = ?`, nu.Id).Scan(&u.Id, &u.Name, &u.Email); err != nil {
		log.Printf("ERROR ReadUser Query() %v", err)
		return nu, err
	}
	return u, nil
}

// UpdateUser implements UserServiceServer
func (s *User) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR UpdateUser RefreshSession() %v", err)
		return nu, err
	}
	defer s.Cluster.Close()
	// update data
	if err := s.Cluster.Session().Query(`UPDATE users SET name = ?, email = ? WHERE id = ?`, nu.Name, nu.Email, nu.Id).Exec(); err != nil {
		log.Printf("ERROR UpdateUser Query() %v", err)
		return nu, err
	}
	return nu, nil
}

// DeleteUser implements UserServiceServer
func (s *User) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR DeleteUser RefreshSession() %v", err)
		return nu, err
	}
	defer s.Cluster.Close()
	// delete data
	if err := s.Cluster.Session().Query(`DELETE FROM users WHERE id = ?`, nu.Id).Exec(); err != nil {
		log.Printf("ERROR DeleteUser Query() %v", err)
		return nu, err
	}
	return &pb.User{}, nil
}

// FilterUsers implements UserServiceServer
func (s *User) FilterUsers(ctx context.Context, f *pb.Filter) (*pb.Users, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR FilterUsers RefreshSession() %v", err)
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
		log.Printf("ERROR FilterUsers Close() %v", err)
		return &pb.Users{}, err
	}

	return users, nil
}

// VerifyUser implements UserServiceServer
func (s *User) VerifyUser(ctx context.Context, crd *pb.Credentials) (*pb.Verification, error) {
	// refresh session
	if err := s.Cluster.RefreshSession(); err != nil {
		log.Printf("ERROR VerifyUser RefreshSession() %v", err)
		return &pb.Verification{IsValid: false}, err
	}
	defer s.Cluster.Close()
	// hash credential password
	challenge := &pb.Credentials{Email: crd.Email, Password: hashPassword(crd.Password)}
	valid := &pb.Credentials{}
	// get data
	if err := s.Cluster.Session().Query(`SELECT email, password FROM users WHERE email = ?`, crd.Email).Scan(&valid.Email, &valid.Password); err != nil {
		log.Printf("ERROR VerifyUser Query() %v", err)
		return &pb.Verification{IsValid: false}, err
	}
	return &pb.Verification{IsValid: challenge == valid}, nil
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
