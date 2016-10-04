package backend

import (
	"crypto/rand"
	"crypto/rsa"

	"bitbucket.org/aukbit/pluto/auth/jwt"
	pba "bitbucket.org/aukbit/pluto/auth/proto"
	"bitbucket.org/aukbit/pluto/client"
	pbu "bitbucket.org/aukbit/pluto/examples/user/proto"
	"golang.org/x/net/context"
)

// Auth struct
type Auth struct {
	Clt client.Client
}

// Authenticate implements authentication
func (a *Auth) Authenticate(ctx context.Context, cre *pba.Credentials) (*pba.Token, error) {
	// make a call to user backend service for credentials verification
	nCred := &pbu.Credentials{Email: cre.Email, Password: cre.Password}
	v, err := a.Clt.Call().(pbu.UserServiceClient).VerifyUser(ctx, nCred)
	if err != nil {
		panic(err)
	}
	if !v.IsValid {
		return &pba.Token{}, nil
	}
	// generate new token
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	token, err := jwt.NewToken(cre.Email, pk)
	if err != nil {
		panic(err)
	}
	return &pba.Token{Jwt: token}, nil
}
