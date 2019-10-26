package backend

import (
	"errors"
	"time"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/auth/jwt"
	pba "github.com/aukbit/pluto/v6/auth/proto"
	"github.com/aukbit/pluto/v6/client"
	pbu "github.com/aukbit/pluto/v6/examples/user/proto"
	"golang.org/x/net/context"
)

var (
	errCredentials            = errors.New("Invalid credentials")
	errClientUserNotAvailable = errors.New("Client user not available")
)

// AuthViews struct
type AuthViews struct{}

// Authenticate implements authentication
func (av *AuthViews) Authenticate(ctx context.Context, cre *pba.Credentials) (*pba.Token, error) {
	// get client user from pluto service from context
	c, ok := pluto.FromContext(ctx).Client("user")
	if !ok {
		return &pba.Token{}, errClientUserNotAvailable
	}
	// dial
	conn, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		return &pba.Token{}, err
	}
	defer conn.Close()
	// make a call to user backend service for credentials verification
	nCred := &pbu.Credentials{Email: cre.Email, Password: cre.Password}
	v, err := c.Stub(conn).(pbu.UserServiceClient).VerifyUser(ctx, nCred)
	if err != nil {
		return &pba.Token{}, err
	}
	if !v.IsValid {
		return &pba.Token{}, errCredentials
	}

	cs := &jwt.ClaimSet{
		Identifier: cre.Email,
		Audience:   "bearer",
		Expiration: 3650,
	}
	token, err := jwt.NewToken(ctx, cs)
	if err != nil {
		return &pba.Token{}, err
	}
	return &pba.Token{Jwt: token}, nil
}

// Verify implements authentication
func (av *AuthViews) Verify(ctx context.Context, t *pba.Token) (*pba.Verification, error) {

	err := jwt.Verify(ctx, t.Jwt)
	if err != nil {
		return &pba.Verification{IsValid: false}, err
	}
	return &pba.Verification{IsValid: true}, nil
}
