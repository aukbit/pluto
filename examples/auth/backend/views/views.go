package backend

import (
	"errors"
	"time"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/auth/jwt"
	pba "github.com/aukbit/pluto/auth/proto"
	"github.com/aukbit/pluto/client"
	pbu "github.com/aukbit/pluto/examples/user/proto"
	"golang.org/x/net/context"
)

var (
	privKeyPath = "./keys/auth.rsa"
	pubKeyPath  = "./keys/auth.rsa.pub"
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
	c, ok := ctx.Value("pluto").(*pluto.Service).Client("user")
	if !ok {
		return &pba.Token{}, errClientUserNotAvailable
	}
	// dial
	i, err := c.Dial(client.Timeout(5 * time.Second))
	if err != nil {
		return &pba.Token{}, err
	}
	defer c.Close()
	// make a call to user backend service for credentials verification
	nCred := &pbu.Credentials{Email: cre.Email, Password: cre.Password}
	v, err := i.(pbu.UserServiceClient).VerifyUser(ctx, nCred)
	if err != nil {
		return &pba.Token{}, err
	}
	if !v.IsValid {
		return &pba.Token{}, errCredentials
	}
	pk, err := jwt.LoadPrivateKey(privKeyPath)
	if err != nil {
		return &pba.Token{}, err
	}
	token, err := jwt.NewToken(cre.Email, pk)
	if err != nil {
		return &pba.Token{}, err
	}
	return &pba.Token{Jwt: token}, nil
}

// Verify implements authentication
func (av *AuthViews) Verify(ctx context.Context, t *pba.Token) (*pba.Verification, error) {
	pk, err := jwt.LoadPublicKey(pubKeyPath)
	if err != nil {
		return &pba.Verification{IsValid: false}, err
	}
	err = jwt.Verify(t.Jwt, pk)
	if err != nil {
		return &pba.Verification{IsValid: false}, err
	}
	return &pba.Verification{IsValid: true}, nil
}
