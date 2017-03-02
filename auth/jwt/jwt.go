package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aukbit/pluto/auth/jws"
)

var (
	privKeyPath = "./auth.rsa"
	pubKeyPath  = "./auth.rsa.pub"
)

// LoadPublicKey loads a public key from PEM encoded data.
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	if path == "" {
		path = pubKeyPath
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		err = errors.New("Invalid PublicKey format")
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// LoadPrivateKey loads a private key from PEM encoded data.
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		path = privKeyPath
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		err = errors.New("Invalid PrivateKey format")
		return nil, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

// NewToken returns a JWT token signed with the given RSA private key.
func NewToken(identifier string, pk *rsa.PrivateKey) (string, error) {
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
	token, err := jws.Encode(header, payload, pk)
	if err != nil {
		return "", err
	}

	return token, nil
}

// Verify tests whether the provided JWT token's signature was produced by the private key
// associated with the supplied public key.
func Verify(token string, key *rsa.PublicKey) error {
	return jws.Verify(token, key)
}

// BearerAuth returns the token provided in the request's
// Authorization header, if the request uses HTTP Bearer Authentication.
func BearerAuth(r *http.Request) (token string, ok bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	return parseBearerAuth(auth)
}

// parseBearerAuth parses an HTTP Bearer Authentication string.
// "Bearer QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("QWxhZGRpbjpvcGVuIHNlc2FtZQ==", true).
func parseBearerAuth(auth string) (token string, ok bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	return auth[len(prefix):], true
}
