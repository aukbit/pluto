package backend

import (
	"crypto/sha256"
	"encoding/hex"

	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/aukbit/pluto/server/ext"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// UserViews struct
type UserViews struct{}

// CreateUser implements UserServiceServer
func (uv *UserViews) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// generate user id uuid
	newID := uuid.New().String()
	// hash password
	passwordHash := hashPassword(nu.Password)
	// persist data
	if err := session.Query(`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		newID, nu.Name, nu.Email, passwordHash).Exec(); err != nil {
		return &pb.User{}, err
	}
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: newID}, nil
}

// ReadUser implements UserServiceServer
func (uv *UserViews) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// user object
	u := &pb.User{}
	// get data
	if err := session.Query(`SELECT id, name, email FROM users WHERE id = ?`, nu.Id).Scan(&u.Id, &u.Name, &u.Email); err != nil {
		return nu, err
	}
	return u, nil
}

// UpdateUser implements UserServiceServer
func (uv *UserViews) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// update data
	if err := session.Query(`UPDATE users SET name = ?, email = ? WHERE id = ?`, nu.Name, nu.Email, nu.Id).Exec(); err != nil {
		return nu, err
	}
	return nu, nil
}

// DeleteUser implements UserServiceServer
func (uv *UserViews) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// delete data
	if err := session.Query(`DELETE FROM users WHERE id = ?`, nu.Id).Exec(); err != nil {
		return nu, err
	}
	return &pb.User{}, nil
}

// FilterUsers implements UserServiceServer
func (uv *UserViews) FilterUsers(ctx context.Context, f *pb.Filter) (*pb.Users, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// filter users
	iter := session.Query(`SELECT id, name, email FROM users WHERE name = ? ALLOW FILTERING;`, f.Name).Iter()
	//
	users := &pb.Users{}
	u := &pb.User{}
	for iter.Scan(&u.Id, &u.Name, &u.Email) {
		users.Data = append(users.Data, u)
	}
	if err := iter.Close(); err != nil {
		return &pb.Users{}, err
	}

	return users, nil
}

// VerifyUser implements UserServiceServer
func (uv *UserViews) VerifyUser(ctx context.Context, crd *pb.Credentials) (*pb.Verification, error) {
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// hash credential password
	challenge := &pb.Credentials{Email: crd.Email, Password: hashPassword(crd.Password)}
	valid := &pb.Credentials{}
	// get data
	if err := session.Query(`SELECT email, password FROM users WHERE email = ?`, crd.Email).Scan(&valid.Email, &valid.Password); err != nil {
		return &pb.Verification{IsValid: false}, err
	}
	return &pb.Verification{IsValid: challenge == valid}, nil
}

// StreamUsers implements UserServiceServer
func (uv *UserViews) StreamUsers(in *pb.Filter, stream pb.UserService_StreamUsersServer) error {
	ctx := stream.Context()
	// get db session from context
	session := ext.FromContextAny(ctx, "cassandra").(*gocql.Session)
	// filter users
	iter := session.Query(`SELECT id, name, email FROM users WHERE name = ? ALLOW FILTERING;`, in.Name).Iter()
	defer iter.Close()

	u := &pb.User{}
	for iter.Scan(&u.Id, &u.Name, &u.Email) {
		// stream
		if err := stream.Send(u); err != nil {
			return err
		}
	}
	return nil
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
