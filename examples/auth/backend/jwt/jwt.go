package backend

import (
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"log"

	"bitbucket.org/aukbit/pluto/jws"
)

const (
	privKeyPath = "/auth.rsa"
	pubKeyPath  = "/auth.rsa.pub"
)

func getPublicKey() ([]byte, error) {
	return ioutil.ReadFile(pubKeyPath)
}

func getPrivateKey() ([]byte, error) {
	return ioutil.ReadFile(pubKeyPath)
}

func newToken(identifier string) string {
	header := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}
	payload := &jws.ClaimSet{
		Iss: identifier,
		Aud: "",
		Exp: 3610,
		Iat: 10,
	}
	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	token, err := jws.Encode(header, payload, pk)
	if err != nil {
		log.Fatal(err)
	}

	return token
}
