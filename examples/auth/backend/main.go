package main

import (
	"log"

	"github.com/aukbit/pluto/v6/auth/jwt"
	"github.com/aukbit/pluto/v6/examples/auth/backend/service"
)

var (
	privKeyPath = "./keys/auth.rsa"
	pubKeyPath  = "./keys/auth.rsa.pub"
)

func main() {
	prv, err := jwt.LoadPrivateKey(privKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	pub, err := jwt.LoadPublicKey(pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	// run backend service
	if err := backend.Run(pub, prv); err != nil {
		log.Fatal(err)
	}
}
