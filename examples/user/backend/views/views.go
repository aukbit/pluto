package backend

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/aukbit/pluto"
	pb "github.com/aukbit/pluto/examples/user/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// UserViews struct
type UserViews struct{}

// CreateUser implements UserServiceServer
func (uv *UserViews) CreateUser(ctx context.Context, nu *pb.NewUser) (*pb.User, error) {
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return &pb.User{}, err
	}
	defer db.Close()
	// generate user id uuid
	newID := uuid.New().String()
	// hash password
	passwordHash := hashPassword(nu.Password)
	// persist data
	if err := db.Session().Query(`INSERT INTO users (id, name, email, password) VALUES (?, ?, ?, ?)`,
		newID, nu.Name, nu.Email, passwordHash).Exec(); err != nil {
		return &pb.User{}, err
	}
	return &pb.User{Name: nu.Name, Email: nu.Email, Id: newID}, nil
}

// ReadUser implements UserServiceServer
func (uv *UserViews) ReadUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return nu, err
	}
	defer db.Close()
	// user object
	u := &pb.User{}
	// get data
	if err := db.Session().Query(`SELECT id, name, email FROM users WHERE id = ?`, nu.Id).Scan(&u.Id, &u.Name, &u.Email); err != nil {
		return nu, err
	}
	return u, nil
}

// UpdateUser implements UserServiceServer
func (uv *UserViews) UpdateUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return nu, err
	}
	defer db.Close()
	// update data
	if err := db.Session().Query(`UPDATE users SET name = ?, email = ? WHERE id = ?`, nu.Name, nu.Email, nu.Id).Exec(); err != nil {
		return nu, err
	}
	return nu, nil
}

// DeleteUser implements UserServiceServer
func (uv *UserViews) DeleteUser(ctx context.Context, nu *pb.User) (*pb.User, error) {
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return nu, err
	}
	defer db.Close()
	// delete data
	if err := db.Session().Query(`DELETE FROM users WHERE id = ?`, nu.Id).Exec(); err != nil {
		return nu, err
	}
	return &pb.User{}, nil
}

// FilterUsers implements UserServiceServer
func (uv *UserViews) FilterUsers(ctx context.Context, f *pb.Filter) (*pb.Users, error) {
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return &pb.Users{}, err
	}
	defer db.Close()
	// filter users
	iter := db.Session().Query(`SELECT id, name, email FROM users WHERE name = ? ALLOW FILTERING;`, f.Name).Iter()

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
	// get datastore from pluto service from context
	db := ctx.Value("pluto").(pluto.Service).Config().Datastore
	// refresh session
	if err := db.RefreshSession(); err != nil {
		return &pb.Verification{IsValid: false}, err
	}
	defer db.Close()
	// hash credential password
	challenge := &pb.Credentials{Email: crd.Email, Password: hashPassword(crd.Password)}
	valid := &pb.Credentials{}
	// get data
	if err := db.Session().Query(`SELECT email, password FROM users WHERE email = ?`, crd.Email).Scan(&valid.Email, &valid.Password); err != nil {
		return &pb.Verification{IsValid: false}, err
	}
	return &pb.Verification{IsValid: challenge == valid}, nil
}

func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}
