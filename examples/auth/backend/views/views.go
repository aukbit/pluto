package backend

import (
	"errors"
	"log"

	"bitbucket.org/aukbit/pluto/auth/jwt"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/client"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	privKeyPath = "./keys/auth.rsa"
	pubKeyPath  = "./keys/auth.rsa.pub"
)

var (
	errCredentials = errors.New("Invalid credentials")
)

// Auth struct
type Auth struct {
	Clt client.Client
}

// Authenticate implements authentication
func (a *Auth) Authenticate(ctx context.Context, cre *pba.Credentials) (*pba.Token, error) {
	// context event
	md, _ := metadata.FromContext(ctx)
	log.Printf("Authenticate %v", md["event"])
	// make a call to user backend service for credentials verification
	nCred := &pbu.Credentials{Email: cre.Email, Password: cre.Password}
	v, err := a.Clt.Call().(pbu.UserServiceClient).VerifyUser(ctx, nCred)
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
func (a *Auth) Verify(ctx context.Context, t *pba.Token) (*pba.Verification, error) {
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
